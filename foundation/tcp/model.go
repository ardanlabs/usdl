package tcp

import "bytes"

// Set of event types.
const (
	EvtAccept = iota + 1
	EvtJoin
	EvtRead
	EvtRemove
	EvtDrop
	EvtGroom
	EvtStop
)

// Set of event sub types.
const (
	TypError = iota + 1
	TypInfo
)

var eventTypes = map[int]string{
	EvtAccept: "accept",
	EvtJoin:   "join",
	EvtRead:   "read",
	EvtRemove: "remove",
	EvtDrop:   "drop",
	EvtGroom:  "groom",
	EvtStop:   "stop",
}

var eventSubTypes = map[int]string{
	TypError: "error",
	TypInfo:  "info",
}

// =============================================================================

// Error provides support for multi client operations that might error.
type Errors []error

// Error implments the error interface for CltError.
func (ers Errors) Error() string {
	var b bytes.Buffer

	for _, err := range ers {
		b.WriteString(err.Error())
		b.WriteString("\n")
	}

	return b.String()
}
