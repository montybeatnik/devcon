package devcon

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSHClient holds the elements to setup an SSH client
type SSHClient struct {
	ip, port  string
	clientCfg *ssh.ClientConfig
}

type option func() (func(*SSHClient), error)

// failedOption - helper to expose error from option builder
func failedOption(err error) option {
	return func() (func(*SSHClient), error) {
		return nil, err
	}
}

// properOption - helper to expose valid setter from option builder
func properOption(setter func(*SSHClient)) option {
	return func() (func(*SSHClient), error) {
		return setter, nil
	}
}

// Password sets SSHClient's password.
func SetPort(port string) option {
	return properOption(func(c *SSHClient) {
		c.port = port
	})
}

// Password sets SSHClient's password.
func SetPassword(pw string) option {
	return properOption(func(c *SSHClient) {
		authMethod := []ssh.AuthMethod{
			ssh.Password(pw),
		}
		c.clientCfg.Auth = authMethod
	})
}

// PrivateKey sets SSHClient's private key.
func SetPrivateKey(keyfile string) option {
	privKeyData, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return failedOption(err)
	}
	privkey, err := ssh.ParsePrivateKey(privKeyData)
	if err != nil {
		return failedOption(err)
	}
	return properOption(func(c *SSHClient) {
		authMethod := []ssh.AuthMethod{
			ssh.PublicKeys(privkey),
		}
		c.clientCfg.Auth = authMethod
	})
}

// Timeout sets SSHClient's timeout value.
func SetTimeout(seconds time.Duration) option {
	return properOption(func(c *SSHClient) {
		c.clientCfg.Timeout = seconds
	})
}

func SetHostKeyCallback(knownHostsFile string) option {
	hostKeyCallback, err := knownhosts.New(knownHostsFile)
	if err != nil {
		return failedOption(err)
	}
	return properOption(func(c *SSHClient) {
		c.clientCfg.HostKeyCallback = hostKeyCallback
	})
}

// setup - helper to set optional dependencies.
func (c *SSHClient) setup(opts ...option) error {
	// today, setup will be repeated for every service with options :(
	if c == nil {
		return nil
	}
	for _, option := range opts {
		if option == nil {
			continue
		}
		setter, err := option()
		if err != nil {
			return err
		}
		if setter != nil {
			setter(c)
		}
	}
	return nil
}

// NewClient is a factory function that takes in SSH parameters
// and returns a new client
func NewClient(un, ip string, opts ...option) (*SSHClient, error) {
	// establish the SSH config from the crytpo package and associate it to
	// the clientCfg field.
	defaultPort := "22"
	client := &SSHClient{
		ip:   ip,
		port: defaultPort,
		clientCfg: &ssh.ClientConfig{
			User:            un,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(time.Second * 5),
		},
	}
	if err := client.setup(opts...); err != nil {
		return client, err
	}
	return client, nil
}

// Run takes in a command and attempts to establishe a remote session
// and run the command.
func (c *SSHClient) Run(cmd string) (string, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.ip, c.port), c.clientCfg)
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
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", c.ip, c.port), c.clientCfg)
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
