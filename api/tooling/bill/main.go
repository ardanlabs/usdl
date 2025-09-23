package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
)

func main() {
	data := []string{
		"I'm happy.",
		"She exercises every morning.",
		"His dog barks loudly.",
		"My school starts at 8:00.",
		"We always eat dinner together.",
		"They take the bus to work.",
		"He doesn't like vegetables.",
		"I don't want anything to drink.",
		"This little black dress isn't expensive.",
		"Those kids don't speak English.",
	}

	for _, d := range data {
		l := len(d)

		lBuf := make([]byte, 2)
		binary.LittleEndian.PutUint16(lBuf, uint16(l))
		fmt.Println("CALLING KEVIN:", l, d)

		var buf []byte
		buf = append(buf, lBuf...)

		v := rand.Intn(10)
		if v == 0 {
			v = 1
		}

		buf = append(buf, []byte(d[:v])...)
		buf = kevin(buf)

		buf = append(buf, []byte(d[v:])...)
		kevin(buf)
	}
}

var isLength = true
var chunkLen int
var written int

func kevin(buf []byte) []byte {
	for {
		if len(buf) == 0 {
			return nil
		}

		if isLength {
			uVal := binary.LittleEndian.Uint16(buf[:2])
			chunkLen = int(uVal)
			buf = buf[2:]
			isLength = false

			fmt.Println("NEW LINE:", chunkLen)
		}

		if len(buf) == 0 {
			return nil
		}

		switch {
		case len(buf) < (chunkLen - written):
			fmt.Printf("WRITE %d BYTES TO DISK 1\n", len(buf))
			fmt.Println(string(buf))
			buf = buf[len(buf):]
			written += len(buf)

		case len(buf) == (chunkLen - written):
			fmt.Printf("WRITE %d BYTES TO DISK 2\n", len(buf))
			fmt.Println(string(buf))
			written = 0
			buf = buf[len(buf):]
			isLength = true

		case len(buf) > (chunkLen - written):
			diff := chunkLen - written
			fmt.Printf("WRITE %d BYTES TO DISK 3\n", diff)
			fmt.Println(string(buf[:diff]))
			buf = buf[diff:]
			written = 0
			isLength = true
		}
	}
}
