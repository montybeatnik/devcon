# Devcon

## Overview
A module to ssh into devices and run commands. It supports running a single command in a remote session or running several commands in an interactive session. 

## Juniper Package
There is also a Juniper package that wraps around the base devcon package that has some logic built in to get structured data back from the devices. It also allows you to perform a "dry run" by logging in, applying a config, printing a diff, and then rolling it back before logging out of the device. 

## Authentication Methods
This package supports the following authentication methods:
- via username/password
- via an SSH key file


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
	devIP := "10.0.0.60"
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		devIP,
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
	devIP := "10.0.0.60"
	client := devcon.NewClient(
		os.Getenv("SSH_USER"),
		devIP,
		devcon.PrivateKey(keyFile),
	)
	out, err := client.Run("show version")
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
	fmt.Println(out)
}

// Simple concurrency example with wait group

```

## TODO
- [x] Add support for public key auth.
- [x] Add Juniper package 
  - [x] Config Differ
  - [x] Apply Configs
  - [x] First Operational Command Example

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