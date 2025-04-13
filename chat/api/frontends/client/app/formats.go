package app

import "fmt"

func errorMessage(format string, a ...any) Message {
	return Message{
		Name:    "system",
		Content: fmt.Appendf(nil, format, a...),
	}
}
