package devcon

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SSHClient holds the elements to setup an SSH client
type SSHClient struct {
	ipAndPort string
	clientCfg *ssh.ClientConfig
}

// NewClient is a factory function that takes in SSH parameters
// and returns a new client
func NewClient(un, ip string, opts ...option) *SSHClient {
	// establish the SSH config from the crytpo package and associate it to
	// the clientCfg field.
	defaultPort := "22"
	ipAndPort := fmt.Sprintf("%v:%v", ip, defaultPort)
	client := &SSHClient{
		ipAndPort: ipAndPort,
		clientCfg: &ssh.ClientConfig{
			User:            un,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(time.Second * 5),
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

type option func(*SSHClient)

// Option sets the options specified.
func (c *SSHClient) Option(opts ...option) {
	for _, opt := range opts {
		opt(c)
	}
}

// Password sets SSHClient's password.
func Password(pw string) option {
	return func(c *SSHClient) {
		authMethod := []ssh.AuthMethod{
			ssh.Password(pw),
		}
		c.clientCfg.Auth = authMethod
	}
}

// Run takes in a command and attempts to establishe a remote session
// and run the command.
func (c *SSHClient) Run(cmd string) (string, error) {
	client, err := ssh.Dial("tcp", c.ipAndPort, c.clientCfg)
	if err != nil {
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
	client, err := ssh.Dial("tcp", c.ipAndPort, c.clientCfg)
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
