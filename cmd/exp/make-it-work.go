package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func main() {
	// user := os.Getenv("SSH_USER")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	keyFile := filepath.Join(homeDir, ".ssh/id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	privKeyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatal(err)
	}
	privkey, err := ssh.ParsePrivateKey(privKeyData)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &ssh.ClientConfig{
		User: "rolodev",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privkey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", "10.0.0.60:22", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	sess, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	sess.Stdout = os.Stdout
	// sess.Stdin = os.Stdin
	_ = sess.Run("show version")
	sess.Close()
}
