package tools

import (
	"bytes"
	"strings"
)

func IsStringInSlice(needle string, haystack []string) (found bool) {
	for _, v := range haystack {
		if strings.EqualFold(needle, v) {
			found = true
			break
		}
	}
	return
}

// a byte buffer that implements the Close() method
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	// this is just memory, so all we need to do is return
	return nil
}
