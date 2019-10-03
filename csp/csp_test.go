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

func TestS32_SQUASH(t *testing.T) {
	characters := "Hello,* ** *CSP."
	expected := "Hello,* ↑ *CSP."

	west, east := make(chan rune), make(chan rune)
	go csp.S32_SQUASH(west, east)
	go func() {
		for _, c := range characters {
			west <- c
		}
		close(west)
	}()
	received := make([]rune, 0, len(expected))
	for r := range east {
		received = append(received, r)
	}
	if string(received) != expected {
		t.Fatalf("%v: expected: %v, got: %v", t.Name(), expected, string(received))
	}
}

func TestS33_DISASSEMBLE(t *testing.T) {
	cardfiles := [][]rune{
		[]rune("Hello,CSP"),
		[]rune("Hello,CSP2"),
		[]rune("Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,"),
	}
	expected := [][]rune{
		[]rune("Hello,CSP"),
		[]rune("Hello,CSP2"),
		[]rune("Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,Hellooo,"),
	}

	cardfile, X := make(chan []rune), make(chan rune)
	go csp.S33_DISASSEMBLE(cardfile, X)
	go func() {
		for _, cf := range cardfiles {
			cardfile <- cf
		}
		close(cardfile)
	}()

	i := 0
	received := []rune{}
	for c := range X {
		if c == ' ' {
			if string(received) != string(expected[i]) {
				t.Fatalf("%v: expected: '%v', got: '%v'", t.Name(), string(expected[i]), string(received))
			}
			i++
			received = []rune{}
			continue
		}
		received = append(received, c)
	}
}
