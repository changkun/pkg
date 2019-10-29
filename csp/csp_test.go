package csp_test

import (
	"reflect"
	"sync"
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
		got := <-re
		if !reflect.DeepEqual(tt.want, got) {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), tt.want, got)
		}
	}
}

func TestS42_Factorial(t *testing.T) {
	tests := []struct {
		limit int
		user  int
		want  int
	}{
		{
			limit: 5,
			user:  0,
			want:  1,
		},
		{
			limit: 5,
			user:  1,
			want:  1,
		},
		{
			limit: 5,
			user:  2,
			want:  2,
		},
		{
			limit: 5,
			user:  3,
			want:  6,
		},
		{
			limit: 5,
			user:  4,
			want:  24,
		},
		{
			limit: 5,
			user:  5,
			want:  120,
		},
		{
			limit: 5,
			user:  6,
			want:  720,
		},
		{
			limit: 5,
			user:  7,
			want:  2520,
		},
	}

	for _, tt := range tests {
		fac := make([]chan int, tt.limit+1)
		for i := 0; i <= tt.limit; i++ {
			fac[i] = make(chan int)
		}
		csp.S42_Factorial(fac, tt.limit)

		fac[0] <- tt.user

		got := <-fac[0]
		if got != tt.want {
			t.Fatalf("%v: expected: %v, got: %v", t.Name(), tt.want, got)
		}
	}
}

func TestS43_SmallSetOfIntegers(t *testing.T) {
	set := csp.NewS43_SmallSetOfIntegers()

	for i := 0; i < 100; i++ {
		done, has := make(chan bool), make(chan bool)
		go set.Insert(i, done)
		if <-done {
			go set.Has(i, has)
			if <-has {
				continue
			}
			t.Fatalf("%v: Has not found, value: %v", t.Name(), i)
		}
		t.Fatalf("%v: Insert undone, value: %v", t.Name(), i)
	}

	for i := 100; i < 200; i++ {
		done := make(chan bool)
		go set.Insert(i, done)
		if <-done {
			t.Fatalf("%v: Insert should not done, but done, value: %v", t.Name(), i)
		}
	}

	for i := 100; i < 200; i++ {
		has := make(chan bool)
		go set.Has(i, has)
		if <-has {
			t.Fatalf("%v: Has should not found, but found, value: %v", t.Name(), i)
		}
	}
}

func TestS44_Scan(t *testing.T) {
	set := csp.NewS43_SmallSetOfIntegers()

	for i := 0; i < 100; i++ {
		done, has := make(chan bool), make(chan bool)
		go set.Insert(i, done)
		if <-done {
			go set.Has(i, has)
			if <-has {
				continue
			}
			t.Fatalf("%v: Has not found, value: %v", t.Name(), i)
		}
		t.Fatalf("%v: Insert undone, value: %v", t.Name(), i)
	}

	readch := make(chan int)
	go set.S44_Scan(readch)

	i := 0
	for v := range readch {
		if v != i {
			t.Fatalf("%v: scan op read inconsistency, expect %v, got %v", t.Name(), i, v)
		}
		i++
	}
}

func TestS45_RecursiveSmallSetOfIntegers(t *testing.T) {
	has, insert, least := csp.S45_S46_NewRecursiveSmallSetOfIntegers()
	for i := 0; i < 100; i++ {
		check := make(chan bool)
		has <- csp.S45_Has{V: i, Response: check}
		if ok := <-check; ok {
			t.Fatalf("%v: expected not has, got value: %v", t.Name(), i)
		}
	}

	for i := 0; i < 100; i++ {
		insert <- i

		check := make(chan bool)
		has <- csp.S45_Has{V: i, Response: check}
		if ok := <-check; !ok {
			t.Fatalf("%v: expected inserted, got false, value: %v", t.Name(), i)
		}

		response := make(chan csp.S46_Least)
		least <- response
		v := <-response
		if !v.NoneLeft {
			t.Fatalf("%v: expecting left, got NoneLeft, value: %v", t.Name(), i)
		}

		if v.Least != i {
			t.Fatalf("%v: expected least %v, got false, value: %v", t.Name(), i, v.Least)
		}
	}

	close(has)
	close(insert)
	close(least)
}

func TestS51_BoundedBuffer(t *testing.T) {
	pro, con := csp.S51_BoundedBuffer()
	go func() {
		for i := 0; i < 100; i++ {
			pro <- i
		}
	}()
	for i := 0; i < 100; i++ {
		v := <-con
		if i != v {
			t.Fatalf("%v: expected %v, got %v", t.Name(), i, v)
		}
	}
}

func TestS52_IntegerSemaphore(t *testing.T) {
	sem := csp.NewS52_IntegerSemaphore()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 100; i++ {
			sem.P()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 100; i++ {
			sem.V()
		}
		wg.Done()
	}()
	wg.Wait()

	sem.Close()
}

func TestS53_DiningPhilosophers(t *testing.T) {
	csp.S53_DiningPhilosophers()
}

func TestS61_TheSieveOfEratosthenes(t *testing.T) {
	csp.S61_TheSieveOfEratosthenes(100)
}
func TestS62_MatrixMultiplication(t *testing.T) {
	A := [][]int{
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		[]int{7, 8, 9},
	}
	m := csp.S62_NewMatrix(A)
	println(1)
	IN := [][]int{
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
	}

	if !reflect.DeepEqual(IN, m.S62_Multiply(IN)) {
		t.Fatalf("matrix multiplication result is incorrect")
	}
}
