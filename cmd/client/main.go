package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/montybeatnik/devcon"
)

var (
	homeLabIP = "10.0.0.60"
	mac       = "10.0.0.80"
	prodCisco = "172.30.83.76"
	localhost = "127.0.0.1"
)

func main() {
	run()
}

type SoftwareVersionReply struct {
	SoftwareInformation struct {
		Text               string `xml:",chardata"`
		HostName           string `xml:"host-name"`
		ProductModel       string `xml:"product-model"`
		ProductName        string `xml:"product-name"`
		Jsr                string `xml:"jsr"`
		PackageInformation struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name"`
			Comment string `xml:"comment"`
		} `xml:"package-information"`
	} `xml:"software-information"`
}

func run() {
	un := os.Getenv("USER")
	// pw := "password"
	pw := os.Getenv("PASSWORD")
	client := devcon.NewClient(un,
		homeLabIP,
		devcon.WithPassword(pw),
	)
	output, err := client.Run("show version | display xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
		os.Exit(42)
	}
	var swver SoftwareVersionReply
	xml.Unmarshal([]byte(output), &swver)
	fmt.Println(swver.SoftwareInformation.HostName, swver.SoftwareInformation.ProductModel)
	// cmds := []string{
	// 	"show ip interface brief",
	// 	"exit",
	// }
	// interfaceOutput, err := client.RunAll(cmds...)
	// if err != nil {
	// 	log.Println("command failed", err)
	// 	os.Exit(42)
	// }
	// fmt.Println(interfaceOutput)
}
