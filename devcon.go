package devcon

import (
	"time"

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
