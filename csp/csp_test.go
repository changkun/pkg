package csp_test

import (
	"reflect"
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

func TestS32_SQUASH_EX(t *testing.T) {
	characters := "Hello,* ** *CSP.***"
	expected := "Hello,* ↑ *CSP.↑*"

	west, east := make(chan rune), make(chan rune)
	go csp.S32_SQUASH_EX(west, east)
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

func TestS34_ASSEMBLE(t *testing.T) {
	tests := []struct {
		stream string
		want   []string
	}{
		{
			stream: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			want: []string{
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
				"67890                                                                                                                        ",
			},
		},
		{
			stream: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234512345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
			want: []string{
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
			},
		},
		{
			stream: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123451234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234",
			want: []string{
				"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
				"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234 ",
			},
		},
		{
			stream: " ",
			want: []string{
				"                                                                                                                             ",
			},
		},
		{
			stream: "",
			want:   []string{},
		},
	}

	for _, tt := range tests {
		ttt := tt
		X, lineprinter := make(chan rune), make(chan string)
		go csp.S34_ASSEMBLE(X, lineprinter)
		go func() {
			for _, c := range ttt.stream {
				X <- c
			}
			close(X)
		}()

		received := []string{}
		for c := range lineprinter {
			received = append(received, string(c))
		}
		if !reflect.DeepEqual(tt.want, received) {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), ttt.want, received)
		}
	}
}

func TestS36_Reformat(t *testing.T) {
	tests := []struct {
		cardfile [][]rune
		want     []string
	}{
		{
			cardfile: [][]rune{
				[]rune("1234567890123456789012345678901234567890123456789012345678901234567890"),
				[]rune("1234567890123456789012345678901234567890123456789012345678901234567890"),
			},
			want: []string{
				"1234567890123456789012345678901234567890123456789012345678901234567890 123456789012345678901234567890123456789012345678901234",
				"5678901234567890                                                                                                             ",
			},
		},
	}

	for _, tt := range tests {
		ttt := tt
		cardfile, lineprinter := make(chan []rune), make(chan string)
		go csp.S35_Reformat(cardfile, lineprinter)
		go func() {
			for _, c := range ttt.cardfile {
				cardfile <- c
			}
			close(cardfile)
		}()

		received := []string{}
		for c := range lineprinter {
			received = append(received, string(c))
		}
		if !reflect.DeepEqual(tt.want, received) {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), ttt.want, received)
		}
	}
}

func TestS36_ConwayProblem(t *testing.T) {
	tests := []struct {
		cardfile [][]rune
		want     []string
	}{
		{
			cardfile: [][]rune{
				[]rune("Hello,* ** *CSP.***"),
			},
			want: []string{
				"Hello,* ↑ *CSP.↑*                                                                                                            ",
			},
		},
	}

	for _, tt := range tests {
		ttt := tt
		cardfile, lineprinter := make(chan []rune), make(chan string)
		go csp.S36_ConwayProblem(cardfile, lineprinter)
		go func() {
			for _, c := range ttt.cardfile {
				cardfile <- c
			}
			close(cardfile)
		}()

		received := []string{}
		for c := range lineprinter {
			received = append(received, string(c))
		}
		if !reflect.DeepEqual(tt.want, received) {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), ttt.want, received)
		}
	}
}

func TestS41_DivisionWithRemainder(t *testing.T) {
	tests := []struct {
		input csp.S41_In
		want  csp.S41_Out
	}{
		{
			input: csp.S41_In{10, 5},
			want:  csp.S41_Out{2, 0},
		},
		{
			input: csp.S41_In{3, 2},
			want:  csp.S41_Out{1, 1},
		},
		{
			input: csp.S41_In{10, 3},
			want:  csp.S41_Out{3, 1},
		},
	}

	for _, tt := range tests {
		ch := make(chan csp.S41_In)
		re := make(chan csp.S41_Out)
		go csp.S41_DivisionWithRemainder(ch, re)
		ch <- tt.input
		ret := <-re
		if !reflect.DeepEqual(tt.want, ret) {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), tt.want, ret)
		}
	}
}
