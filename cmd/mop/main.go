package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/montybeatnik/devcon"
	"github.com/montybeatnik/devcon/vendors/juniper"
)

func main() {
	for hn, devcfg := range updateConfigs {
		log.Printf("applying config on %v", hn)
		for ip, cfgfile := range devcfg {
			cfg := getConfig(cfgfile)
			log.Printf("here is the cfg: %v\n", cfg)
			dev := getDevice(ip)
			output, err := dev.Diff(cfg)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(output)
		}
	}
}

var updateConfigs = map[string]map[string]string{
	"lab-r1": {"10.0.0.86": "docs/configs/lab-r1.conf"},
	// "lab-r2": {"10.0.0.23": "docs/configs/lab-r2.conf"},
	// "lab-r3": {"10.0.0.212": "docs/configs/lab-r3.conf"},
	// "lab-r4": {"10.0.0.150": "docs/configs/lab-r4.conf"},
}

func getConfig(fn string) []string {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var cfg []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cfg = append(cfg, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return cfg
}

func getDevice(ip string) *juniper.JuniperClient {
	khfp := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	jnprClient := juniper.NewJuniperClient(
		os.Getenv("SSH_USER"),
		ip,
		devcon.WithPassword(os.Getenv("SSH_PASSWORD")),
		devcon.WithTimeout(time.Second*1),
		devcon.WithHostKeyCallback(khfp),
	)
	return jnprClient
	// intTerse, err := jnprClient.InterfacesTerse()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(intTerse)
}
