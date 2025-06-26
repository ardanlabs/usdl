package client

import (
	"bytes"
	"fmt"
)

func StitchMessages(msgs [][]byte) string {
	var b bytes.Buffer
	for _, msg := range msgs {
		b.Write(msg)
	}

	return b.String()
}

func errorMessage(format string, a ...any) Message {
	return Message{
		Name:    "system",
		Content: [][]byte{fmt.Appendf(nil, format, a...)},
	}
}
