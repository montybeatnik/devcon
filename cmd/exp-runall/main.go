package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	user := os.Getenv("SSH_USER")
	pass := os.Getenv("SSH_PASSWORD")
	ipAndPort := "10.0.0.60:22"
	cmds := []string{"show version", "exit"}
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// var output bytes.Buffer
	client, err := ssh.Dial("tcp", ipAndPort, cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Create sesssion
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	defer session.Close()
	// assign reader and writer
	// Store the session output to an io.Reader
	session.Stdout = os.Stdout
	// StdinPipe for commands
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	// Start remote shell
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}
	// send the commands
	for _, cmd := range cmds {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
	// buf, err := io.ReadAll(stdout)
	// if err != nil {
	// 	log.Fatal(errors.Wrap(err, "reader to byte slice failed"))
	// }
}
