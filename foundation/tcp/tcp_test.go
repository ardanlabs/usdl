package tcp_test

import (
	"bufio"
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
		cfg := tcp.Config{
			NetType: "tcp4",
			Addr:    ":0",

			ConnHandler: tcpConnHandler{},
			ReqHandler:  tcpReqHandler{},
			RespHandler: tcpRespHandler{},
		}

		// Create a new TCP value.
		u, err := tcp.New("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Start accepting client data.
		if err := u.Start(); err != nil {
			t.Fatal("\tShould be able to start the TCP listener.", "X", err)
		}
		t.Log("\tShould be able to start the TCP listener.", "OK")

		defer u.Stop()

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
	}
}

// Test tcp.Addr works correctly.
func TestTCPAddr(t *testing.T) {
	t.Log("Given the need to listen on any open port and know that bound address.")
	{
		// Create a configuration.
		cfg := tcp.Config{
			NetType: "tcp4",
			Addr:    ":0", // Defer port assignment to OS.

			ConnHandler: tcpConnHandler{},
			ReqHandler:  tcpReqHandler{},
			RespHandler: tcpRespHandler{},
		}

		// Create a new TCP value.
		u, err := tcp.New("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Addr should be nil before Start.
		if addr := u.Addr(); addr != nil {
			t.Fatalf("\tAddr() should be nil before Start; Addr() = %q. %s", addr, "X")
		}
		t.Log("\tAddr() should be nil before Start.", "OK")

		// Start accepting client data.
		if err := u.Start(); err != nil {
			t.Fatal("\tShould be able to start the TCP listener.", "X", err)
		}
		defer u.Stop()

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
	}
}

// TestDropConnections tests we can drop connections when configured.
func TestDropConnections(t *testing.T) {
	t.Log("Given the need to drop TCP connections.")
	{
		// Create a configuration.
		cfg := tcp.Config{
			NetType:     "tcp4",
			Addr:        ":0",
			ConnHandler: tcpConnHandler{},
			ReqHandler:  tcpReqHandler{},
			RespHandler: tcpRespHandler{},
		}

		// Create a new TCP value.
		u, err := tcp.New("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Set the drop connection flag to true.
		t.Log("\tSet the drop connections flag to TRUE.", "OK")
		u.DropConnections(true)

		// Start accepting client data.
		if err := u.Start(); err != nil {
			t.Fatal("\tShould be able to start the TCP listener.", "X", err)
		}
		t.Log("\tShould be able to start the TCP listener.", "OK")

		defer u.Stop()

		// Let's connect to the host:port.
		conn, err := net.Dial("tcp4", u.Addr().String())
		if err != nil {
			t.Fatal("\tShould be able to dial a new TCP connection.", "X", err)
		}
		t.Log("\tShould be able to dial a new TCP connection.", "OK")

		// An attempt to read should result in an EOF.
		b := make([]byte, 1)
		if _, err = conn.Read(b); err == nil {
			t.Fatal("\tShould not be able to read the response from the connection.", "X", err)
		}
		t.Log("\tShould not be able to read the response from the connection.", "OK")
	}
}

// TestRateLimit tests we can drop connections when they come in too fast.
func TestRateLimit(t *testing.T) {
	const ratelimit = 1 * time.Second

	t.Log("Given the need to drop TCP connections.")
	{
		// Create a configuration.
		cfg := tcp.Config{
			NetType:       "tcp4",
			Addr:          ":0",
			ConnHandler:   tcpConnHandler{},
			ReqHandler:    tcpReqHandler{},
			RespHandler:   tcpRespHandler{},
			ConnRateLimit: func() time.Duration { return ratelimit },
		}

		// Create a new TCP value.
		u, err := tcp.New("TEST", cfg)
		if err != nil {
			t.Fatal("\tShould be able to create a new TCP listener.", "X", err)
		}
		t.Log("\tShould be able to create a new TCP listener.", "OK")

		// Start accepting client data.
		if err := u.Start(); err != nil {
			t.Fatal("\tShould be able to start the TCP listener.", "X", err)
		}
		t.Log("\tShould be able to start the TCP listener.", "OK")

		defer u.Stop()

		newconn := func() (*bufio.Writer, *bufio.Reader, net.Conn, error) {
			// Let's connect to the host:port.
			conn, err := net.Dial("tcp4", u.Addr().String())
			if err != nil {
				return nil, nil, nil, err
			}
			return bufio.NewWriter(conn), bufio.NewReader(conn), conn, nil
		}

		// Make a successful connection
		successfulTest := func(Context interface{}) {
			w, r, c, err := newconn()
			if err != nil {
				t.Fatal("\tShould be able to dial a new TCP connection.", Context, "X", err)
			}
			t.Log("\tShould be able to dial a new TCP connection.", Context, "OK")

			defer c.Close()

			if _, err := w.WriteString("Hello\n"); err != nil {
				t.Fatal("\tShould be able to send data to the connection.", Context, "X", err)
			}
			t.Log("\tShould be able to send data to the connection.", Context, "OK")

			if err := w.Flush(); err != nil {
				t.Fatal("\tShould be able to flush the writer.", Context, "X", err)
			}
			t.Log("\tShould be able to flush the writer.", Context, "OK")

			// Let's read the response.
			response, err := r.ReadString('\n')
			if err != nil {
				t.Fatal("\tShould be able to read the response from the connection.", Context, "X", err)
			}
			t.Log("\tShould be able to read the response from the connection.", Context, "OK")

			t.Log(response)
		}

		successfulTest("PRE-LIMIT")

		// The next 100 connections should fail (assuming it's all under rateLimit amount of time).
		for i := 0; i < 100; i++ {
			//  Apparently, even though the connection should not exist, we are still allowed
			//  to connect to the remote socket and write to it.  The error is exhibited
			//  only when it's time to perform a read on that connection.
			w, r, c, err := newconn()
			if err != nil {
				t.Fatal("\tShould be able to dial a non-first TCP connection.", "X", err)
			}
			t.Log("\tShould be able to dial a non-first TCP connection", c.LocalAddr(), "OK")

			defer c.Close()

			if _, err := w.WriteString("Hello\n"); err != nil {
				t.Fatal("\tShould be able to send data to the connection.", "X", err)
			}
			t.Log("\tShould be able to send data to the connection.", "OK")

			if err := w.Flush(); err != nil {
				t.Fatal("\tShould be able to flush the writer.", "X", err)
			}
			t.Log("\tShould be able to flush the writer.", "OK")

			// Let's read the response.
			_, err = r.ReadString('\n')
			if err == nil {
				t.Fatal("\tShould have to read from the connection.", "X")
			}
			t.Log("\tShould have to read from the connection", "OK", err)
		}

		// Sleep for rateLimit to perform another successful test.
		time.Sleep(ratelimit)
		successfulTest("POST-LIMIT")

		// NOTE If you call another 'successfulTest' here, we will fail because we expect
		// the test to fail due to the limit.
	}
}
