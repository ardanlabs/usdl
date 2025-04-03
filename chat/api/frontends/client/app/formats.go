package app

import "fmt"

func formatMessage(name string, msg []byte) []byte {
	return fmt.Appendf(nil, "%s: %s", name, string(msg))
}
