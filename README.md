# Devcon
A simple package to ssh into devices.

## Usage
```go
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
```

## TODO
- [x] Pull the client bits out of RunCommand and pull that logic into the Factory function

## Profile
- go test -v -run Run -cpuprofile cpu.prof -memprofile mem.prof -bench .

## Testing
### Unit tests
#### All tests
go test -v
#### A specific test
go test -v -run RunCommand
#### Specific tests matching a pattern
go test -v -run Run
### Benchmarks
go test -run RunCommand -bench=.