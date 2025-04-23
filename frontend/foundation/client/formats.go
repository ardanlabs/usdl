package client

import "fmt"

// func formatMessage(name string, msg []byte) []byte {
// 	return fmt.Appendf(nil, "%s: %s", name, string(msg))
// }

func errorMessage(format string, a ...any) Message {
	return Message{
		Name:    "system",
		Content: fmt.Appendf(nil, format, a...),
	}
}
