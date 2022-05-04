package devcon

import (
	"os"
	"strings"
	"testing"
)

func getCredsFromEnv() (string, string) {
	return os.Getenv("SSH_USER"), os.Getenv("SSH_PASSWORD")
}

func TestNewClient(t *testing.T) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		t.Fatal("env variables not set")
	}
	client := NewClient(un, pw, "127.0.0.1")
	if client.ip == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.ip, os.Getenv("SSH_USER"))
	}
	t.Logf("\nuser: %q\nIP: %q", client.clientCfg.User, client.ip)
}

func TestRunCommand(t *testing.T) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		t.Fatal("env variables not set")
	}
	client := NewClient(un, pw, "172.28.48.13")
	output, err := client.RunCommand("show version")
	if err != nil {
		t.Fatal(err)
	}
	delimeter := strings.Repeat("#", 60)
	t.Logf("\n%v\n%v\n%v", delimeter, output, delimeter)
}
