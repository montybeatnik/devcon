package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/montybeatnik/devcon"
)

func main() {
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		"10.0.0.60",
		devcon.Password(os.Getenv("SSH_PASSWORD")),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	keyFile := filepath.Join(homeDir, ".ssh/id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	client := devcon.NewClient(
		"rolodev",
		"10.0.0.60",
		devcon.PrivateKey(keyFile),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}

// func main() {
// 	// user := os.Getenv("SSH_USER")
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	keyFile := filepath.Join(homeDir, ".ssh/id_rsa")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	privKeyData, err := ioutil.ReadFile(keyFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	privkey, err := ssh.ParsePrivateKey(privKeyData)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	cfg := &ssh.ClientConfig{
// 		User: "rolodev",
// 		Auth: []ssh.AuthMethod{
// 			ssh.PublicKeys(privkey),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}
// 	conn, err := ssh.Dial("tcp", "10.0.0.60:22", cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	sess, err := conn.NewSession()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	sess.Stdout = os.Stdout
// 	sess.Stdin = os.Stdin
// 	_ = sess.Run("show version")
// 	sess.Close()
// }
