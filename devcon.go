package devcon

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	ErrAuthFailure = errors.New("failed authentication")
	ErrTimeout     = errors.New("timeout")
)

// SSHClient holds the elements to setup an SSH client
type SSHClient struct {
	target, port string
	clientCfg    *ssh.ClientConfig
}

type option func(*SSHClient)

// WithPort sets SSHClient's listening port.
func WithPort(port string) option {
	return func(c *SSHClient) {
		c.port = port
	}
}

// WithPassword sets SSHClient's password.
func WithPassword(pw string) option {
	return func(c *SSHClient) {
		authMethod := []ssh.AuthMethod{
			ssh.Password(pw),
		}
		c.clientCfg.Auth = authMethod
	}
}

// WithPrivateKey sets SSHClient's private key.
func WithPrivateKey(keyfile string) option {
	privKeyData, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatal(err)
	}
	privkey, err := ssh.ParsePrivateKey(privKeyData)
	if err != nil {
		log.Fatal(err)
	}
	return func(c *SSHClient) {
		authMethod := []ssh.AuthMethod{
			ssh.PublicKeys(privkey),
		}
		c.clientCfg.Auth = authMethod
	}
}

// WithTimeout sets SSHClient's timeout value.
func WithTimeout(seconds time.Duration) option {
	return func(c *SSHClient) {
		c.clientCfg.Timeout = seconds
	}
}

// WithHostKeyCallback sets the SSHClient's initializes the
// client with an allow list of known trusted hosts.
func WithHostKeyCallback(knownHostsFile string) option {
	return func(c *SSHClient) {
		hostKeyCallback, err := knownhosts.New(knownHostsFile)
		if err != nil {
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		}
		c.clientCfg.HostKeyCallback = hostKeyCallback
	}
}

// NewClient is a factory function that takes in SSH parameters
// and returns a new client
func NewClient(user, target string, opts ...option) *SSHClient {
	// establish the SSH config from the crytpo package and associate it to
	// the clientCfg field.
	defaultPort := "22"
	client := &SSHClient{
		target: target,
		port:   defaultPort,
		clientCfg: &ssh.ClientConfig{
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(time.Second * 5),
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Run takes in a command and attempts to establishe a remote session
// and run the command.
func (c *SSHClient) Run(cmd string) (string, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.target, c.port), c.clientCfg)
	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			return "", ErrAuthFailure
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			return "", ErrTimeout
		}
		return "", errors.Wrap(err, "dail failed")
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "NewSession failed")
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), errors.Wrap(err, "Run failed")
	}
	return string(output), err
}

// RunAll takes in one or more commands. It establishes a remote session with the target IP
// and attempts to run all of the commands supplied. You must remember to exit as this
// method does establish an interactive session.
func (c *SSHClient) RunAll(cmds ...string) (string, error) {
	output, err := c.executeMany(cmds...)
	if err != nil {
		return "", err
	}
	return output.String(), nil
}

// assignStdInAndOut takes in a pointer to an ssh.Session and returns a Reader,
// WriteCloser, and an error
func assignStdInAndOut(sess *ssh.Session) (io.Reader, io.WriteCloser, error) {
	// Store the session output to an io.Reader
	sshOut, err := sess.StdoutPipe()
	if err != nil {
		var sIn io.WriteCloser
		return sshOut, sIn, fmt.Errorf("failed to get stdOut: %v", err)
	}
	// StdinPipe for commands
	stdIn, err := sess.StdinPipe()
	if err != nil {
		return sshOut, stdIn, fmt.Errorf("failed to get stdIn: %v", err)
	}
	return sshOut, stdIn, nil
}

// executeMany sets up an interactive session with the target device
func (c *SSHClient) executeMany(cmds ...string) (bytes.Buffer, error) {
	var output bytes.Buffer
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.target, c.port), c.clientCfg)
	if err != nil {
		return output, err
	}
	// Create sesssion
	session, err := client.NewSession()
	if err != nil {
		return output, err
	}
	defer client.Close()
	defer session.Close()
	// assign reader and writer
	stdOut, stdin, err := assignStdInAndOut(session)
	if err != nil {
		return output, err
	}
	// Start remote shell
	err = session.Shell()
	if err != nil {
		return output, err
	}
	// send the commands
	for _, cmd := range cmds {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return output, fmt.Errorf(fmt.Sprintf("Failed to get cmd output: %v", err))
		}
	}
	err = session.Wait()
	if err != nil {
		return output, fmt.Errorf(fmt.Sprintf("Failed to exit: %v", err))
	}
	buf, err := io.ReadAll(stdOut)
	if err != nil {
		return output, errors.Wrap(err, "reader to byte slice failed")
	}
	return *bytes.NewBuffer(buf), nil
}
