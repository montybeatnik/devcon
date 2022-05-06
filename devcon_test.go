package devcon

import (
	"os"
	"strings"
	"testing"
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
	client, err := NewClient(un, pw, labSRX)
	if err != nil {
		t.Fatal(err)
	}
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
	client, err := NewClient(un, pw, labSRX)
	if err != nil {
		t.Fatal(err)
	}
	output, err := client.RunCommand("show version")
	if err != nil {
		t.Fatal(err)
	}
	delimeter := strings.Repeat("#", 60)
	t.Logf("\n%v\n%v\n%v", delimeter, output, delimeter)
}

func BenchmarkRunCommand(b *testing.B) {
	un, pw := getCredsFromEnv()
	if un == "" || pw == "" {
		b.Error("env variables not set")
	}
	client, err := NewClient(un, pw, labSRX)
	if err != nil {
		b.Fatal(err)
	}
	// run the RunCommand method b.N times
	for n := 0; n < b.N; n++ {
		_, err := client.RunCommand("show version")
		if err != nil {
			b.Log(err)
		}
	}
}
