package utils

import (
	"fmt"
	"io"
)

func CheckWriter(w io.Writer) error {
	const data = "test"
	n, err := w.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("write: short write: %d", n)
	}
	if w, ok := w.(io.Closer); ok {
		err := w.Close()
		if err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}
	return nil
}
