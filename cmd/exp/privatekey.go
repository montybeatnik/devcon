package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/montybeatnik/devcon"
)

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
		os.Getenv("SSH_USER"),
		"10.0.0.86",
		devcon.WithPrivateKey(keyFile),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}
