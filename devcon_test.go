package devcon

import (
	"io"
	"net"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
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
	client := NewClient(un, labSRX, SetPassword(pw))
	if client.ipAndPort == "" {
		t.Fatalf("\ngot: %q | wanted %q\n", client.ipAndPort, os.Getenv("SSH_USER"))
	}
	t.Logf("\nuser: %q\nIP: %q", client.clientCfg.User, client.ipAndPort)
}

// func TestRunCommand(t *testing.T) {
// 	un, pw := getCredsFromEnv()
// 	if un == "" || pw == "" {
// 		t.Fatal("env variables not set")
// 	}
// 	client := NewClient(un, labSRX, SetHostKeyCallback(pw))
// 	output, err := client.Run("show version")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	delimeter := strings.Repeat("#", 60)
// 	t.Logf("\n%v\n%v\n%v", delimeter, output, delimeter)
// }

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

// netPipe is analogous to net.Pipe, but it uses a real net.Conn, and
// therefore is buffered (net.Pipe deadlocks if both sides start with
// a write.)
func netPipe() (net.Conn, net.Conn, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		listener, err = net.Listen("tcp", "[::1]:0")
		if err != nil {
			return nil, nil, err
		}
	}
	defer listener.Close()
	c1, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		return nil, nil, err
	}

	c2, err := listener.Accept()
	if err != nil {
		c1.Close()
		return nil, nil, err
	}

	return c1, c2, nil
}

type serverType func(ssh.Channel, <-chan *ssh.Request, *testing.T)

// dial constructs a new test server and returns a *ClientConn.
func dial(handler serverType, t *testing.T) *ssh.Client {
	c1, c2, err := netPipe()
	if err != nil {
		t.Fatalf("netPipe: %v", err)
	}

	go func() {
		defer c1.Close()
		conf := ssh.ServerConfig{
			NoClientAuth: true,
		}
		conf.AddHostKey(testSigners["rsa"])

		conn, chans, reqs, err := ssh.NewServerConn(c1, &conf)
		t.Log("here")
		if err != nil {
			t.Logf("Unable to handshake: %v", err)
		}
		go ssh.DiscardRequests(reqs)

		for newCh := range chans {
			if newCh.ChannelType() != "session" {
				newCh.Reject(1, "unknown channel type")
				continue
			}

			ch, inReqs, err := newCh.Accept()
			if err != nil {
				t.Errorf("Accept: %v", err)
				continue
			}
			go func() {
				handler(ch, inReqs, t)
			}()
		}
		if err := conn.Wait(); err != io.EOF {
			t.Logf("server exit reason: %v", err)
		}
	}()

	config := &ssh.ClientConfig{
		User:            "testuser",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, chans, reqs, err := ssh.NewClientConn(c2, "", config)
	if err != nil {
		t.Fatalf("unable to dial remote side: %v", err)
	}

	return ssh.NewClient(conn, chans, reqs)
}

func TestSessionCombinedOutput(t *testing.T) {
	conn := dial(fixedOutputHandler, t)
	defer conn.Close()
	// session, err := conn.NewSession()
	// if err != nil {
	// 	t.Fatalf("Unable to request new session: %v", err)
	// }
	// defer session.Close()

	// buf, err := session.CombinedOutput("") // cmd is ignored by fixedOutputHandler
	// if err != nil {
	// 	t.Error("Remote command did not exit cleanly:", err)
	// }
	// fmt.Println(buf)
	// const stdout = "this-is-stdout."
	// const stderr = "this-is-stderr."
	// g := string(buf)
	// if g != stdout+stderr && g != stderr+stdout {
	// 	t.Error("Remote command did not return expected string:")
	// 	t.Logf("want %q, or %q", stdout+stderr, stderr+stdout)
	// 	t.Logf("got  %q", g)
	// }
}

// Ignores the command, writes fixed strings to stderr and stdout.
// Strings are "this-is-stdout." and "this-is-stderr.".
func fixedOutputHandler(ch ssh.Channel, in <-chan *ssh.Request, t *testing.T) {
	defer ch.Close()
	_, err := ch.Read(nil)

	req, ok := <-in
	if !ok {
		t.Fatalf("error: expected channel request, got: %#v", err)
		return
	}

	// ignore request, always send some text
	req.Reply(true, nil)

	_, err = io.WriteString(ch, "this-is-stdout.")
	if err != nil {
		t.Fatalf("error writing on server: %v", err)
	}
	_, err = io.WriteString(ch.Stderr(), "this-is-stderr.")
	if err != nil {
		t.Fatalf("error writing on server: %v", err)
	}
	sendStatus(0, ch, t)
}

type exitStatusMsg struct {
	Status uint32
}

func sendStatus(status uint32, ch ssh.Channel, t *testing.T) {
	msg := exitStatusMsg{
		Status: status,
	}
	if _, err := ch.SendRequest("exit-status", false, ssh.Marshal(&msg)); err != nil {
		t.Errorf("unable to send status: %v", err)
	}
}
