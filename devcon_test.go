package devcon

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(os.Getenv("SSH_USER"), os.Getenv("SSH_PASSWORD"), "localhost")
	t.Logf("IP = %q", client.ip)
}
