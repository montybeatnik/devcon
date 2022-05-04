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
	ip           string
	port         string
	clientCfg    ssh.ClientConfig
	session      ssh.Session
	clientClose  func() error
	sessionClose func() error
	// output
}

// NewClient is a factory function that takes in SSH parameters
// and returns a new client
func NewClient(un, pw, ip string, args ...string) (*sshClient, error) {
	// establish the SSH config from the crytpo package and associate it to
	// the clientCfg field.
	clientCfg := &ssh.ClientConfig{
		User: un,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(time.Second * 5),
	}
	defaultPort := "22"
	ipAndPort := fmt.Sprintf(ip + ":" + defaultPort)
	client, err := ssh.Dial("tcp", ipAndPort, clientCfg)
	if err != nil {
		return nil, errors.Wrap(err, "dail failed")
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "NewSession failed")
	}
	return &sshClient{
		ip:           ip,
		clientCfg:    *clientCfg,
		session:      *session,
		clientClose:  client.Close,
		sessionClose: session.Close,
	}, nil
}

// RunCommand takes in a command and attempts to establishe a session
// and run the command. Should the session or command fail, a meaningful
// error is returned. If the command succeeds, the output and a nil error
// is returned.
func (c *sshClient) RunCommand(cmd string) (string, error) {
	// ensure the client and session are closed
	defer c.clientClose()
	defer c.sessionClose()
	var outBuf bytes.Buffer
	c.session.Stdout = &outBuf
	err := c.session.Run(cmd)
	if err != nil {
		return outBuf.String(), errors.Wrap(err, "Run failed")
	}
	return outBuf.String(), err
}
