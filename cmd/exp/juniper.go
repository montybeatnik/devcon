package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
	"github.com/montybeatnik/devcon/vendors/juniper"
)

func main() {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	jnprClient := juniper.NewJuniperClient(
		os.Getenv("SSH_USER"),
		"10.0.0.212",
		devcon.WithPassword(os.Getenv("SSH_PASSWORD")),
		devcon.WithTimeout(time.Second*1),
		devcon.WithHostKeyCallback(khfp),
	)
	intTerse, err := jnprClient.InterfacesTerse()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(intTerse)
}
