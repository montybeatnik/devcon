package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
)

func main() {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	// knownHostsFile, err := os.Open(khfp)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer knownHostsFile.Close()
	// log.Println(knownHostsFile)
	client, err := devcon.NewClient(
		os.Getenv("SSH_USER"),
		"10.0.0.60",
		// devcon.SetPassword(os.Getenv("SSH_PASSWORD")),
		devcon.SetPrivateKey("blah"),
		devcon.SetTimeout(time.Second*1),
		devcon.SetHostKeyCallback(khfp),
	)
	if err != nil {
		log.Fatal("couldn't spin up client: ", err)
	}
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}

	fmt.Println(out)
}
