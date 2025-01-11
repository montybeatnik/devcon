package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
)

/*
# lab nodes

- vsrx1 - 10.0.0.86
- vsrx2 - 10.0.0.23
- vsrx3 - 10.0.0.212
*/

func main() {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	// knownHostsFile, err := os.Open(khfp)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer knownHostsFile.Close()
	// log.Println(knownHostsFile)
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		"10.0.0.212",
		devcon.WithPassword(os.Getenv("SSH_PASSWORD")),
		devcon.WithTimeout(time.Second*1),
		devcon.WithHostKeyCallback(khfp),
	)

	// ####################################
	// ########## interface terse #########
	// ####################################
	// cmd := "show interfaces terse | display xml"
	// out, err := client.Run(cmd)
	// if err != nil {
	// 	log.Fatalf("command failed: %v", err)
	// }

	// var intTerse juniper.InterfaceTerse
	// if err := xml.Unmarshal([]byte(out), &intTerse); err != nil {
	// 	log.Println(err)
	// }

	// for _, ifd := range intTerse.InterfaceInformation.PhysicalInterface {
	// 	fmt.Println(ifd.Name)
	// }

	// for _, ifd := range intTerse.InterfaceInformation.PhysicalInterface {
	// 	for _, ifl := range ifd.LogicalInterface {
	// 		fmt.Println(ifl.Name)
	// 		for _, addr := range ifl.AddressFamily {
	// 			fmt.Printf("%v | ", addr.AddressFamilyName)
	// 			for _, a := range addr.InterfaceAddress {
	// 				fmt.Println(a.IfaLocal.Text)
	// 			}
	// 		}
	// 	}
	// }
	// ####################################

	cfg := []string{
		"configure",
		"set protocols lldp interface all",
		"show | compare",
		"commit and-quit comment ADDING_LLDP",
		"exit",
	}
	output, err := client.RunAll(cfg...)
	if err != nil {
		log.Printf("cfg failed", err)
	}
	fmt.Println(output)
}
