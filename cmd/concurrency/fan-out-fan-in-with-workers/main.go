package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/montybeatnik/devcon"
	"github.com/pkg/profile"
)

var (
	devicesCSV = `us-iowac-rtr08,172.28.48.13
bo3cl-rtr02,172.30.83.68
us-dn1vw-rtr01,172.30.115.7
us-dn1vw-rtr02,172.30.114.38
us-iowac-rtr01,172.30.117.26`

	pw = os.Getenv("SSH_PASSWORD")
	un = os.Getenv("SSH_USER")
)

var workers int

func init() {
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of workers (defaults to # of logical CPUs).")
}

func main() {
	defer profile.Start(
		profile.TraceProfile,
		profile.ProfilePath("."),
	).Stop()
	start := time.Now()
	devices := strings.Split(devicesCSV, "\n")
	// The done channel will be shared by the entire pipeline
	// so that when it's closed it serves as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	defer close(done)

	in := gen(done, devices...)

	// fan-out
	var chans []<-chan auditOp
	for i := 0; i < workers; i++ {
		chans = append(chans, audit(done, in))
	}
	for s := range merge(done, chans...) {
		if s.auditErr != nil {
			fmt.Println("[FAILED] - ", s.hostname, s.auditErr)
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

func gen(done <-chan struct{}, lines ...string) <-chan auditOp {
	out := make(chan auditOp, len(lines))
	go func() {
		defer close(out)
		for _, d := range lines {
			select {
			case <-done:
				return
			default:
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
		}
	}()
	return out
}

func audit(done <-chan struct{}, in <-chan auditOp) <-chan auditOp {
	out := make(chan auditOp)
	go func() {
		defer close(out)
		for audit := range in {
			select {
			default:
				client := devcon.NewClient(un, audit.ip, devcon.SetPassword(pw))
				start := time.Now()
				output, err := client.Run("show version | i IOS XE ")
				audit.duration = time.Since(start)
				if err != nil {
					audit.auditErr = err
				} else {
					audit.output = output
				}
				out <- audit
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, chans ...<-chan auditOp) <-chan auditOp {
	out := make(chan auditOp)
	wg := sync.WaitGroup{}
	wg.Add(len(chans))

	for _, sc := range chans {
		go func(sc <-chan auditOp) {
			defer wg.Done()
			for audit := range sc {
				select {
				case out <- audit:
				case <-done:
					return
				}
			}
		}(sc)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
