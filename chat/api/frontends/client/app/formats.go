package app

import "fmt"

func formatMessage(name string, msg []byte) []byte {
	return []byte(fmt.Sprintf("%s: %s", name, string(msg)))
}
