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
	// var wg *sync.WaitGroup
	// wg.Add(len(updateConfigs))
	for hn, devcfg := range updateConfigs {
		log.Printf("applying config on %v", hn)
		for ip, cfgfile := range devcfg {
			// go func(cfgfile string, ip string) {
			// 	defer wg.Done()
			cfg := getConfig(cfgfile)
			dev := getDevice(ip)
			// output, err := dev.Diff(cfg)
			output, err := dev.ApplyConfig(cfg)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(output)

			// }(cfgfile, ip)
		}
	}
	// wg.Wait()
}

var updateConfigs = map[string]map[string]string{
	"lab-r1": {"10.0.0.86": "docs/configs/lab-r1.conf"},
	"lab-r2": {"10.0.0.23": "docs/configs/lab-r2.conf"},
	"lab-r3": {"10.0.0.212": "docs/configs/lab-r3.conf"},
	"lab-r4": {"10.0.0.150": "docs/configs/lab-r4.conf"},
	"lab-r5": {"10.0.0.87": "docs/configs/lab-r5.conf"},
	"lab-r6": {"10.0.0.24": "docs/configs/lab-r6.conf"},
	"lab-r7": {"10.0.0.213": "docs/configs/lab-r7.conf"},
	"lab-r8": {"10.0.0.149": "docs/configs/lab-r8.conf"},
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
