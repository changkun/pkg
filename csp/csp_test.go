package csp_test

import (
	"testing"

	"github.com/changkun/gobase/csp"
)

func TestS31_COPY(t *testing.T) {
	characters := "Hello, CSP."

	west, east := make(chan rune), make(chan rune)
	go csp.S31_COPY(west, east)
	go func() {
		for _, c := range characters {
			west <- c
		}
		close(west)
	}()

	received := make([]rune, 0, len(characters))
	for r := range east {
		received = append(received, r)
	}
	if string(received) != characters {
		t.Fatalf("%v: expected: %v, got: %v", t.Name(), characters, string(received))
	}
}
