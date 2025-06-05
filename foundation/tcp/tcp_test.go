package tcp_test

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

// TestTCP provide a test of listening for a connection and
// echoing the data back.
func TestTCP(t *testing.T) {
	t.Log("Given the need to listen and process TCP data.")
	{
		// Create a configuration.
		cfg := tcp.ServerConfig{
			NetType:  "tcp4",
			Addr:     ":0",
			Handlers: tcpHandlers{},
		}

		// Create a new TCP value.
		u, err := tcp.NewServer("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Start accepting client data.
		go func() {
			if err := u.Listen(); err != nil {
				t.Error("\tShould be able to start the TCP listener.", "X", err)
			}
			t.Log("\tShould be able to start the TCP listener.", "OK")
		}()

		// Let's connect back and send a TCP package
		conn, err := net.Dial("tcp4", u.Addr().String())
		if err != nil {
			t.Fatal("\tShould be able to dial a new TCP connection.", "X", err)
		}
		t.Log("\tShould be able to dial a new TCP connection.", "OK")

		// Setup a bufio reader to extract the response.
		bufReader := bufio.NewReader(conn)
		bufWriter := bufio.NewWriter(conn)

		// Send some know data to the tcp listener.
		if _, err := bufWriter.WriteString("Hello\n"); err != nil {
			t.Fatal("\tShould be able to send data to the connection.", "X", err)
		}
		t.Log("\tShould be able to send data to the connection.", "OK")

		bufWriter.Flush()

		// Let's read the response.
		response, err := bufReader.ReadString('\n')
		if err != nil {
			t.Fatal("\tShould be able to read the response from the connection.", "X", err)
		}
		t.Log("\tShould be able to read the response from the connection.", "OK")

		if response == "GOT IT\n" {
			t.Log("\tShould receive the string \"GOT IT\".", "OK")
		} else {
			t.Error("\tShould receive the string \"GOT IT\".", "X", response)
		}

		d := atomic.LoadInt64(&dur)
		duration := time.Duration(d)

		if duration <= 2*time.Second {
			t.Log("\tShould be less that 2 seconds.", "OK")
		} else {
			t.Error("\tShould be less that 2 seconds.", "X", duration)
		}

		u.Shutdown(context.Background())
	}
}

// Test tcp.Addr works correctly.
func TestTCPAddr(t *testing.T) {
	t.Log("Given the need to listen on any open port and know that bound address.")
	{
		// Create a configuration.
		cfg := tcp.ServerConfig{
			NetType:  "tcp4",
			Addr:     ":0", // Defer port assignment to OS.
			Handlers: tcpHandlers{},
		}

		// Create a new TCP value.
		u, err := tcp.NewServer("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Addr should be nil before Start.
		if addr := u.Addr(); addr != nil {
			t.Fatalf("\tAddr() should be nil before Start; Addr() = %q. %s", addr, "X")
		}
		t.Log("\tAddr() should be nil before Start.", "OK")

		go func() {
			// Start accepting client data.
			if err := u.Listen(); err != nil {
				t.Error("\tShould be able to start the TCP listener.", "X", err)
			}
		}()

		// Addr should be non-nil after Start.
		addr := u.Addr()
		if addr == nil {
			t.Fatal("\tAddr() should be not be nil after Start.", "X")
		}
		t.Log("\tAddr() should be not be nil after Start.", "OK")

		// The OS should assign a random open port, which shouldn't be 0.
		_, port, err := net.SplitHostPort(addr.String())
		if err != nil {
			t.Fatalf("\tSplitHostPort should not fail. %v. %s", err, "X")
		}
		if port == "0" {
			t.Fatalf("\tAddr port should not be %q. %s", port, "X")
		}
		t.Logf("\tAddr() should be not be 0 after Start (port = %q). %s", port, "OK")

		u.Shutdown(context.Background())
	}
}
