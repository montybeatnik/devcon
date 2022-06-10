package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/montybeatnik/devcon"
)

var (
	devicesCSV = `us-iowac-rtr08,172.28.48.13
bo3cl-rtr02,172.30.83.68
us-dn1vw-rtr01,172.30.115.7
us-dn1vw-rtr02,172.30.114.38
us-iowac-rtr01,172.30.117.26`

	pw = os.Getenv("PASSWORD")
	un = os.Getenv("USER")
)

func main() {
	start := time.Now()
	devices := strings.Split(devicesCSV, "\n")
	scanChan := audit(gen(devices...))
	for s := range scanChan {
		if s.auditErr != nil {
			fmt.Println(s.hostname, s.auditErr)
			continue
		}
		fmt.Println(s.hostname, s.output, s.duration)
	}
	fmt.Println("took:", time.Since(start))
}

type auditOp struct {
	hostname string
	ip       string
	output   string
	auditErr error
	duration time.Duration
}

func gen(lines ...string) <-chan auditOp {
	out := make(chan auditOp, len(lines))
	go func() {
		defer close(out)
		for _, d := range lines {
			hnip := strings.Split(d, ",")
			hn := hnip[0]
			ip := hnip[1]
			out <- auditOp{
				hostname: hn,
				ip:       ip,
				auditErr: nil,
				output:   "",
			}
		}
	}()
	return out
}

func audit(in <-chan auditOp) <-chan auditOp {
	out := make(chan auditOp)
	go func() {
		defer close(out)
		for audit := range in {
			client := devcon.NewClient(un, audit.ip, devcon.SetPassword(pw))
			start := time.Now()
			output, err := client.Run("show version | i IOS ")
			audit.duration = time.Since(start)
			if err != nil {
				audit.auditErr = err
			} else {
				audit.output = output
			}
			out <- audit
		}
	}()
	return out
}
