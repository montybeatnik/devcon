package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		switch s.RawCommand() {
		case "show version":
			io.WriteString(s, fmt.Sprintf("17.1R3\n"))
		case "show interfaces terse":
			io.WriteString(s, fmt.Sprintf("ge-0/0/0"))
		default:
			io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
		}
	})

	log.Fatal(ssh.ListenAndServe("10.0.0.80:22", nil,
		ssh.HostKeyFile("/Users/chrishern/.ssh/id_rsa"),
		ssh.PasswordAuth(ssh.PasswordHandler(func(ctx ssh.Context, password string) bool {
			return password == "password"
		})),
	))
}

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}
