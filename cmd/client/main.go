package main

import (
	"fmt"
	"log"
	"os"

	"github.com/montybeatnik/devcon"
)

func main() {
	un := os.Getenv("SSH_USER")
	pw := os.Getenv("SSH_PASSWORD")
	ip := "10.0.0.60"
	client, err := devcon.NewClient(un, pw, ip)
	if err != nil {
		log.Println("client setup failed", err)
		os.Exit(42)
	}
	output, err := client.RunCommand("show version")
	if err != nil {
		log.Println("command failed", err)
		os.Exit(42)
	}
	fmt.Println(output)
}
