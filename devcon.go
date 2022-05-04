package devcon

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// sshClient holds the elements to setup an SSH client
type sshClient struct {
	ip        string
	port      string
	clientCfg *ssh.ClientConfig
}

// NewClient is a factory function that takes in SSH parameters
// and returns a new client
func NewClient(un, pw, ip string, args ...string) *sshClient {
	return &sshClient{
		ip: ip,
		// establish the SSH config from the crytpo package and associate it to
		// the clientCfg field.
		clientCfg: &ssh.ClientConfig{
			User: un,
			Auth: []ssh.AuthMethod{
				ssh.Password(pw),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(time.Second * 5),
		},
	}
}

// RunCommand takes in a command and attempts to establishe a session
// and run the command. Should the session or command fail, a meaningful
// error is returned. If the command succeeds, the output and a nil error
// is returned.
func (c *sshClient) RunCommand(cmd string) (string, error) {
	var outBuf bytes.Buffer
	c.port = "22"
	// Concat the ip and the port for input to the Dial func.
	ipAndPort := fmt.Sprintf(c.ip + ":" + c.port)
	client, err := ssh.Dial("tcp", ipAndPort, c.clientCfg)
	if err != nil {
		return outBuf.String(), err
	}
	// Ensure the client is closed.
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return outBuf.String(), errors.Wrap(err, "NewSession failed")
	}
	// Ensure the session is closed.
	defer session.Close()
	// session.Stdout is a Writer interface. We're assigning a concrete
	// implmentation to the abstract field.
	session.Stdout = &outBuf
	err = session.Run(cmd)
	if err != nil {
		return outBuf.String(), errors.Wrap(err, "Run failed")
	}
	return outBuf.String(), err
}
