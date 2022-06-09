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

// Example with Password
func main() {
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		"10.0.0.60",
		devcon.Password(os.Getenv("SSH_PASSWORD")),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}

// Example with Private Key
func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	keyFile := filepath.Join(homeDir, ".ssh/id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	client := devcon.NewClient(
		"rolodev",
		"10.0.0.60",
		devcon.PrivateKey(keyFile),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}
```

## TODO
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