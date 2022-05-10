package devcon

import (
	"fmt"
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
func NewClient(un, pw, ip string, args ...string) *SSHClient {
	// establish the SSH config from the crytpo package and associate it to
	// the clientCfg field.
	defaultPort := "22"
	ipAndPort := fmt.Sprintf(ip + ":" + defaultPort)
	return &SSHClient{
		ipAndPort: ipAndPort,
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
func (c *SSHClient) RunCommand(cmd string) (string, error) {
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
