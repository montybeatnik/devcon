package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/montybeatnik/devcon"
	"github.com/pkg/profile"
)

var prodCisco = "172.30.83.76"

var devicesCSV = `us-iowac-rtr08,172.28.48.13
bo3cl-rtr02,172.30.83.68
us-dn1vw-rtr01,172.30.115.7
us-dn1vw-rtr02,172.30.114.38
us-iowac-rtr01,172.30.117.26`

func main() {
	defer profile.Start(
		profile.TraceProfile,
		profile.ProfilePath("."),
	).Stop()
	start := time.Now()
	pw := os.Getenv("SSH_PASSWORD")
	un := os.Getenv("SSH_USER")
	devices := strings.Split(devicesCSV, "\n")
	for _, d := range devices {
		hnip := strings.Split(d, ",")
		hn := hnip[0]
		ip := hnip[1]
		client := devcon.NewClient(un, ip, devcon.SetPassword(pw))
		out, err := client.Run("show version | i IOS XE")
		if err != nil {
			log.Println(err)
		}
		fmt.Println(hn, out)
	}
	fmt.Println("took:", time.Since(start))
}
