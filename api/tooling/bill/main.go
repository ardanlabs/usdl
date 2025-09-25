package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
)

func main() {
	data := []string{
		"I'm happy. ",
		"She exercises every morning. ",
		"His dog barks loudly. ",
		"My school starts at 8:00. ",
		"We always eat dinner together. ",
		"They take the bus to work. ",
		"He doesn't like vegetables. ",
		"I don't want anything to drink. ",
		"This little black dress isn't expensive. ",
		"Those kids don't speak English.",
	}

	var streamBuf bytes.Buffer
	for _, d := range data {
		l := len(d)
		lenD := make([]byte, 2)
		binary.LittleEndian.PutUint16(lenD, uint16(l))

		streamBuf.Write(lenD)
		streamBuf.WriteString(d)

		fmt.Print(d)
	}

	fmt.Print("\n\n=======================================\n\n")

	stream := streamBuf.Bytes()

	for len(stream) > 0 {
		v := min(len(stream), min(rand.Intn(20), len(stream)))
		kevin2([]byte(stream[:v]))
		stream = stream[v:]
	}
}

var buffer = make([]byte, 1024*32)
var bufLen = 0

func kevin2(buf []byte) {
	const prefix = 2

	// STORE NEW BYTES INTO THE BUFFER
	n := copy(buffer[bufLen:], buf)
	bufLen += n

	// DO WE HAVE ENOUGH BYTES FOR LENGTH?
	if len(buffer) <= 1 {
		return
	}

	// DO WE HAVE ALL THE EXPECTED BYTES
	uVal := binary.LittleEndian.Uint16(buffer[:prefix])
	chunkLen := int(uVal)
	if chunkLen > bufLen-prefix {
		return
	}

	// WRITE BYTES TO DISK
	fmt.Print(string(buffer[prefix : chunkLen+prefix]))

	// MOVE BYTES TO FRONT OF BUFFER
	n = copy(buffer, buffer[prefix+chunkLen:bufLen])
	bufLen = n
}

// func kevin(buf []byte) {
// 	buffer = append(buffer, buf...)

// 	for len(buffer) > 0 {
// 		if isLength {
// 			// WE DON'T HAVE THE FULL 2 BYTES OF LENGTH INFORMATION.
// 			if len(buffer) <= 1 {
// 				break
// 			}

// 			// WE HAVE THE 2 BYTES FOR THE LENGTH.
// 			uVal := binary.LittleEndian.Uint16(buffer[:2])
// 			chunkLen = int(uVal)
// 			buffer = buffer[2:]
// 			isLength = false
// 		}

// 		// IF WE HAVE NO MORE BYTES TO PROCESS RETURN TO GET MORE.
// 		if len(buffer) == 0 {
// 			return
// 		}

// 		switch {
// 		// WE HAVE LESS BYTES IN THE BUFFER THAN WE ARE EXPECTING.
// 		case len(buffer) < (chunkLen - written):
// 			fmt.Print(string(buffer))
// 			written += len(buffer)
// 			buffer = buffer[len(buffer):]

// 		// WE HAVE THE EXACT NUMBER OF BYTES IN THE BUFFER WE NEED.
// 		case len(buffer) == (chunkLen - written):
// 			fmt.Print(string(buffer))
// 			written = 0
// 			buffer = buffer[len(buffer):]
// 			isLength = true

// 		// WE HAVE MORE BYTES THAN WE EXPECT.
// 		case len(buffer) > (chunkLen - written):
// 			diff := chunkLen - written
// 			fmt.Print(string(buffer[:diff]))
// 			written = 0
// 			buffer = buffer[diff:]
// 			isLength = true
// 		}
// 	}
// }
