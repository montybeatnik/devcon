package devcon

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/gliderlabs/ssh"
)

var labSRX = "10.0.0.60"

func getCredsFromEnv() (string, string) {
	return os.Getenv("SSH_USER"), os.Getenv("SSH_PASSWORD")
}

func TestNewClient(t *testing.T) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		t.Fatal("env variables not set")
	}
	ip := "127.0.0.1"
	port := "22"
	client := NewClient(un, ip, SetPassword(pw))
	if client.ip == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.ip, ip)
	}
	if client.port == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.port, port)
	}
}

func serverUp(ctx context.Context) {
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

	for {
		select {
		default:
			// log.Fatal(ssh.ListenAndServe("127.0.0.1:2222", nil,
			ssh.ListenAndServe("127.0.0.1:2222", nil,
				ssh.HostKeyFile("/Users/chrishern/.ssh/id_rsa"),
				ssh.PasswordAuth(ssh.PasswordHandler(func(ctx ssh.Context, password string) bool {
					return password == "password"
				})),
			)
		}
	}
}

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func TestRunCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go serverUp(ctx)
	time.Sleep(time.Second * 1)
	client := NewClient(
		"timmy",
		"127.0.0.1",
		SetPassword("password"),
		SetPort("2222"),
	)
	output, err := client.Run("show version")
	if err != nil {
		t.Fatal(err)
	}
	expected := "17.1R3"
	if !strings.Contains(output, expected) {
		t.Errorf("got: %q, expected: %q", strings.TrimSpace(output), expected)
	}
}

func BenchmarkRunCommand(b *testing.B) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		b.Error("env variables not set")
	}
	client := NewClient(un, labSRX, SetPassword(pw))
	// run the RunCommand method b.N times
	for n := 0; n < b.N; n++ {
		_, err := client.Run("show version")
		if err != nil {
			b.Log(err)
		}
	}
}
