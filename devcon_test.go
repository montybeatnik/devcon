package netgo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	glssh "github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh"
)

var (
	localhost = "localhost"
)

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
	client := NewClient(un, ip, WithPassword(pw))
	if client.target == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.target, ip)
	}
	if client.port == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.port, port)
	}
}

var (
	verCMD           = "show version"
	verRespons       = "17.1R3"
	intfTerseCMD     = "show interfaces terse"
	intfTerseRespons = "ge-0/0/0"
)

func serverUp(ctx context.Context) {
	glssh.Handle(func(s glssh.Session) {
		switch s.RawCommand() {
		case verCMD:
			io.WriteString(s, verRespons)
		case intfTerseCMD:
			io.WriteString(s, intfTerseRespons)
		default:
			io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
		}
	})

	for {
		select {
		default:
			// log.Fatal(ssh.ListenAndServe("127.0.0.1:2222", nil,
			glssh.ListenAndServe("127.0.0.1:2222", nil,
				glssh.HostKeyFile("/Users/chrishern/.ssh/id_rsa"),
				glssh.PasswordAuth(glssh.PasswordHandler(func(ctx glssh.Context, password string) bool {
					return password == "password"
				})),
			)
		}
	}
}

func TestRunCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go serverUp(ctx)
	time.Sleep(time.Second * 1)
	testCases := map[string]struct {
		target   string
		pw       string
		cmd      string
		expected string
		err      error
	}{
		"version": {
			target:   localhost,
			cmd:      verCMD,
			expected: verRespons,
			pw:       "password",
		},
		"pw-fail": {
			target:   localhost,
			cmd:      verCMD,
			expected: "",
			pw:       "haha",
			err:      ErrAuthFailure,
		},
		"interface terse": {
			target:   localhost,
			cmd:      intfTerseCMD,
			expected: intfTerseRespons,
			pw:       "password",
		},
		"timeout": {
			target:   "1.1.1.1",
			cmd:      "",
			expected: "",
			pw:       "password",
			err:      ErrTimeout,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := NewClient(
				"timmy",
				tc.target,
				WithPassword(tc.pw),
				WithPort("2222"),
			)
			output, err := client.Run(tc.cmd)
			if err != nil {
				if err != tc.err {
					t.Fatal(err)
				}
				t.SkipNow()
			}
			if output != tc.expected {
				t.Errorf("got: %q, expected: %q", strings.TrimSpace(output), tc.expected)
			}
		})
	}
}

func TestSetTimeout(t *testing.T) {
	dur := time.Second * 5
	client := NewClient("blah", localhost, WithTimeout(dur))
	if client.clientCfg.Timeout != dur {
		t.Errorf("got: %v, expected: %v", client.clientCfg.Timeout, dur)
	}
}

func TestSetPort(t *testing.T) {
	port := "2222"
	client := NewClient("blah", localhost, WithPort(port))
	if client.port != port {
		t.Errorf("got: %v, expected: %v", client.port, port)
	}
}

func TestSetPrivateKey(t *testing.T) {
	keyfile := filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	client := NewClient("blah", localhost, WithPrivateKey(keyfile))
	for _, am := range client.clientCfg.Auth {
		if am == nil {
			t.Errorf("expected non nil auth method")
		}
	}
}

func TestRunAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go serverUp(ctx)
	time.Sleep(time.Second * 1)
	client := NewClient("chrishern", localhost,
		WithPassword("password"),
		WithPort("2222"),
	)
	_, err := client.RunAll()
	if err != nil {
		if !strings.Contains(err.Error(), "connection refused") {
			t.Error(err)
		}
	}
}

func TestAssignStdInAndOut(t *testing.T) {
	_, _, err := assignStdInAndOut(&ssh.Session{})
	if err != nil {
		t.Error(err)
	}
}

func TestHostKeyCallback(t *testing.T) {
	kh := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	client := NewClient("chrishern", localhost,
		WithHostKeyCallback(kh),
	)
	if client.clientCfg.HostKeyCallback == nil {
		t.Error("host key call ack should not be nil")
	}
}

func BenchmarkRunCommand(b *testing.B) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		b.Error("env variables not set")
	}
	client := NewClient(un, localhost, WithPassword(pw))
	for n := 0; n < b.N; n++ {
		_, err := client.Run("show version")
		if err != nil {
			b.Log(err)
		}
	}
}
