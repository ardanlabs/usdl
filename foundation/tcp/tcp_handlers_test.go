package tcp_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

// tcpConnHandler is required to process data.
type tcpHandlers struct{}

// Bind is called to init to reader and writer.
func (tcpHandlers) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (tcpHandlers) Read(ipAddress string, r io.Reader) ([]byte, int, error) {
	bufReader := r.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	return []byte(line), len(line), nil
}

var dur int64

// Process is used to handle the processing of the message.
func (tcpHandlers) Process(r *tcp.Request, w io.Writer) {
	bufWriter := w.(*bufio.Writer)
	if _, err := bufWriter.WriteString("GOT IT\n"); err != nil {
		fmt.Println("***> SERVER: ERROR SENDING RESPONSE:", err)
		return
	}

	bufWriter.Flush()

	d := int64(time.Since(r.ReadAt))
	atomic.StoreInt64(&dur, d)
}
