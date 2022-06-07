# Devcon
A simple package to ssh into devices.

## Overview
This package only supports authentication via username/password. 

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
	un := os.Getenv("USER")
	pw := os.Getenv("PASSWORD")
	ip := "10.0.0.60"
	client, err := devcon.NewClient(un, pw, ip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create client: %v\n", err)
		os.Exit(42)
	}
	output, err := client.Run("show version")
	if err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
		os.Exit(42)
	}
	fmt.Println(output)
}
```

## TODO
- [] Add Docker container for test server.
- [] Add support for public key auth. 

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