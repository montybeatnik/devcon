package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var (
	mac = "10.0.0.80"
)

func main() {
	processStream := make(chan *os.Process)
	go serverUp(processStream)
	// un := os.Getenv("USER")
	// pw := "password"
	// client := devcon.NewClient(un, pw, mac)

	var wg sync.WaitGroup
	count := 5
	wg.Add(count)
	for i := 0; i < count; i++ {
		time.Sleep(time.Second * 5)
		// go func(wg *sync.WaitGroup) {
		// 	un := os.Getenv("USER")
		// 	// pw := "password"
		// 	client := devcon.NewClient(un, "127.0.0.1", devcon.SetPassword("password"))
		// 	defer wg.Done()
		// 	output, err := client.Run("show version")
		// 	if err != nil {
		// 		fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
		// 		os.Exit(42)
		// 	}
		// 	fmt.Println(output)
		// }(&wg)
	}
	wg.Wait()

	proc := <-processStream
	proc.Kill()

	fmt.Println("exiting")
}

func serverUp(processStream chan *os.Process) {
	d := "/Users/chrishern/github.com/montybeatnik/devcon/cmd/server"
	p := filepath.Join(d, "main.go")
	run := exec.Command("go", "run", p)
	fmt.Println("firing up server...")
	run.Start()
	fmt.Println(run.Process.Pid)

	fmt.Println("server running...")
	processStream <- run.Process
}
