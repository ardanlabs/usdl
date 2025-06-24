package tcp_test

import (
	"bufio"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

// tcpConnHandler is required to process data.
type tcpHandlers struct{}

// Bind is called to init to reader and writer.
func (tcpHandlers) Bind(clt *tcp.Client) {
	clt.Reader = bufio.NewReader(clt.Conn)
}

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (tcpHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	bufReader := clt.Reader.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	return []byte(line), len(line), nil
}

var dur int64

// Process is used to handle the processing of the message.
func (tcpHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	if _, err := clt.Writer.Write([]byte("GOT IT\n")); err != nil {
		fmt.Println("***> SERVER: ERROR SENDING RESPONSE:", err)
		return
	}

	d := int64(time.Since(r.ReadAt))
	atomic.StoreInt64(&dur, d)
}

func (tcpHandlers) Drop(clt *tcp.Client) {
	fmt.Println("***> SERVER: CONNECTION CLOSED")
}
