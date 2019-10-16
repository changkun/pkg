// Package csp implements "Hoare, C. A. R. (1978). Communicating
// sequential processes. Communications of the ACM, 21(8), 666–677.
// https://doi.org/10.1145/359576.359585".
//
// The paper describes a program specification with several commands:
//
// Program structure
//
//   <cmd>               ::= <simple cmd> | <structured cmd>
//   <simple cmd>        ::= <assignment cmd> | <input cmd> | <output cmd>
//   <structured cmd>    ::= <alternative cmd> | <repetitive cmd> | <parallel cmd>
//   <cmd list>          ::= {<declaration>; | <cmd>; } <cmd>
//
// Parallel command
//
//   <parallel cmd>      ::= [<proc>{||<proc>}]
//   <proc>              ::= <proc label> <cmd list>
//   <proc label>        ::= <empty> | <identifier> :: | <identifier>(<label subscript>{,<label subscript>}) ::
//   <label subscript>   ::= <integer const> | <range>
//   <integer constant>  ::= <numeral> | <bound var>
//   <bound var>         ::= <identifier>
//   <range>             ::= <bound variable>:<lower bound>..<upper bound>
//   <lower bound>       ::= <integer const>
//   <upper bound>       ::= <integer const>
//
// Assignment command
//
//   <assignment cmd>    ::= <target var> := <expr>
//   <expr>              ::= <simple expr> | <structured expr>
//   <structured expr>   ::= <constructor> ( <expr list> )
//   <constructor>       ::= <identifier> | <empty>
//   <expr list>         ::= <empty> | <expr> {, <expr> }
//   <target var>        ::= <simple var> | <structured target>
//   <structured target> ::= <constructor> ( <target var list> )
//   <target var list>   ::= <empty> | <target var> {, <target var> }
//
// Input and output command
//
//   <input cmd>         ::= <source> ? <target var>
//   <output cmd>        ::= <destination> ! <expr>
//   <source>            ::= <proc name>
//   <destination>       ::= <proc name>
//   <proc name>         ::= <identifier> | <identifier> ( <subscripts> )
//   <subscripts>        ::= <integer expr> {, <integer expr> }
//
// Repetitive and alternative command
//
//   <repetitive cmd>    ::= * <alternative cmd>
//   <alternative cmd>   ::= [<guarded cmd> { □ <guarded cmd> }]
//   <guarded cmd>       ::= <guard> → <cmd list> | ( <range> {, <range> }) <guard> → <cmd list>
//   <guard>             ::= <guard list> | <guard list>;<input cmd> | <input cmd>
//   <guard list>        ::= <guard elem> {; <guard elem>}
//   <guard elem>        ::= <boolean expr> | <declaration>
//
// Subroutines and Data Representations
//
// A coroutine acting as a subroutine is a process operating
// concurrently with its user process in a prallel command:
//
//   [subr::SUBROUTINE||X::USER]
//
// The SUBROUTINE will contain a repetitive command:
//
//   *[X?(value params) -> ...; X!(result params)]
//
// where ... computes the results from the values input. The subroutine
// will terminate when its user does. The USER will call subroutine by a
// pair of commands:
//
//   subr!(arguments);...;subr?(results)
//
// Any commands between these two will be executed concurrently with the
// subroutine.
//
// You can find the paper proposed solution in the comment of a function.
//
// Author: Changkun Ou <hi@changkun.us>
package csp

// S31_COPY implements Section 3.1 COPY problem:
// "Write a process X to copy characters output by process west to
// process, east."
//
// Solution:
//
//   X :: *[c:character; west?c -> east!c]
func S31_COPY(west, east chan rune) {
	for c := range west {
		east <- c
	}
	close(east)
}

// S32_SQUASH implements Section 3.2 SQUASH problem:
// "Adapt the previous program to replace every pair of consecutive
// asterisks "**" by an upward arrow "↑". Assume that the final
// character input is not an asterisk."
//
// Solution:
//
//   X :: *[c:character; west?c ->
//     [ c != asterisk -> east!c
//      □ c = asterisk -> west?c;
//            [ c != asterisk -> east!asterisk; east!c
//             □ c = asterisk -> east!upward arrow
//     ] ]    ]
func S32_SQUASH(west, east chan rune) {
	for {
		c, ok := <-west
		if !ok {
			break
		}
		if string(c) != "*" {
			east <- c
		}
		if string(c) == "*" {
			c, ok = <-west
			if !ok {
				break
			}
			if string(c) != "*" {
				east <- '*'
				east <- c
			}
			if string(c) == "*" {
				east <- '↑'
			}
		}
	}
	close(east)
}

// S32_SQUASH_EX implements Section 3.2 SQUASH exercise:
// "(2) As an exercise, adapt this process to deal sensibly with input
// which ends with an odd number of asterisks."
//
// Solution:
//
//   X :: *[c:character; west?c ->
//     [ c != asterisk -> east!c
//      □ c = asterisk -> west?c;
//            [ c != asterisk -> east!asterisk; east!c
//             □ c = asterisk -> east!upward arrow
//            ] □ east!asterisk
//     ]   ]
func S32_SQUASH_EX(west, east chan rune) {
	for {
		c, ok := <-west
		if !ok {
			break
		}
		if c != '*' {
			east <- c
		}
		if c == '*' {
			c, ok = <-west
			if !ok {
				east <- '*'
				break
			}
			if c != '*' {
				east <- '*'
				east <- c
			}
			if c == '*' {
				east <- '↑'
			}
		}
	}
	close(east)
}

// S33_DISASSEMBLE implements Section 3.3 DISASSEMBLE problem:
// "to read cards from a cardfile and output to process X the stream of
// characters they contain. An extra space should be inserted at the end
// of each card."
//
// Solution:
//
//   *[cardimage:(1..80)characters; cardfile?cardimage ->
//       i:integer; i := 1;
//       *[i <= 80 -> X!cardimage(i); i := i+1 ]
//       X!space
//   ]
func S33_DISASSEMBLE(cardfile chan []rune, X chan rune) {
	cardimage := make([]rune, 0, 80)
	for tmp := range cardfile {
		if len(tmp) > 80 {
			cardimage = append(cardimage, tmp[:80]...)
		} else {
			cardimage = append(cardimage, tmp[:len(tmp)]...)
		}
		for i := 0; i < len(cardimage); i++ {
			X <- cardimage[i]
		}
		X <- ' '
		cardimage = cardimage[:0]
	}
	close(X)

	// Alternative solution (But wrong):
	// for cardimage := range cardfile {
	// 	for _, c := range cardimage {
	// 		X <- c
	// 	}
	// 	X <- ' '
	// }
	// close(X)
}

// S34_ASSEMBLE implements Section 3.4 ASSEMBLE problem:
// "To read a stream of characters from process X and print them in
// lines of 125 characters on a lineprinter. The last line should be
// completed with spaces if necessary."
//
// Solution:
//
//   lineimage:(1..125)character;
//   i:integer, i:=1;
//   *[c:character; X?c ->
//       lineimage(i) := c;
//       [i <= 124 -> i := i+1
//       □ i = 125 -> lineprinter!lineimage; i:=1
//   ]   ];
//   [ i = 1 -> skip
//   □ i > 1 -> *[i <= 125 -> lineimage(i) := space; i := i+1];
//     lineprinter!lineimage
//   ]
func S34_ASSEMBLE(X chan rune, lineprinter chan string) {
	lineimage := make([]rune, 125)

	i := 0
	for c := range X {
		lineimage[i] = c
		if i <= 124 {
			i++
		}
		if i == 125 {
			lineimage[i-1] = c
			lineprinter <- string(lineimage)
			i = 0
		}
	}
	if i > 0 {
		for i <= 124 {
			lineimage[i] = ' '
			i++
		}
		lineprinter <- string(lineimage)
	}

	close(lineprinter)
	return
}

// S35_Reformat implements Section 3.5 Reformat problem:
// "Read a sequence of cards of 80 characters each, and print the
// characters on a lineprinter at 125 characters per line. Every card
// should be followed by an extra space, and the last line should be
// complete with spaces if necessary."
//
// Solution:
//
//   [west::DISASSEMBLE||X:COPY||east::ASSEMBLE]
func S35_Reformat(cardfile chan []rune, lineprinter chan string) {
	west, east := make(chan rune), make(chan rune)
	go S33_DISASSEMBLE(cardfile, west)
	go S31_COPY(west, east)
	S34_ASSEMBLE(east, lineprinter)
}

// S36_ConwayProblem implements Section 3.6 Conway's Problem:
// "Adapt the above program to replace every pair of consecutive
// asterisk by an upward arrow."
//
// Solution:
//
//   [west::DISASSEMBLE||X::SQUASH||east::ASSEMBLE]
func S36_ConwayProblem(cardfile chan []rune, lineprinter chan string) {
	west, east := make(chan rune), make(chan rune)
	go S33_DISASSEMBLE(cardfile, west)
	go S32_SQUASH_EX(west, east)
	S34_ASSEMBLE(east, lineprinter)
}

type S41_In struct {
	X, Y int
}
type S41_Out struct {
	Quot, Rem int
}

// S41_DivisionWithRemainder implements Section 4.1 Division With
// Remainer.
// "Construct a process to represent a function type subroutine, which
// accepts a positive dividend and divisor, and returns their integer
// quotient and remainder. Efficiency is of no concern."
//
// Solution:
//
//   [DIV::*[x,y:integer; X?(x,y)->
//         quot,rem:integer; quot := 0; rem := x;
//         *[rem >= y -> rem := rem - y; quot := quot + 1;]
//         X!(quot,rem)
//         ]
//   ||X::USER]
func S41_DivisionWithRemainder(in chan S41_In, out chan S41_Out) {
	v := <-in
	x, y := v.X, v.Y

	quot, rem := 0, x
	for rem >= y {
		rem -= y
		quot++
	}
	out <- S41_Out{quot, rem}
}

// S42_Factorial implements Section 4.2 Factorial
// "Compute a factorial by the recursive method, to a given limit."
//
// Solution:
//
//   [fac(i:1..limit)::
//   *[n:integer;fac(i-1)?n ->
//     [n=0 -> fac(i-1)!1
//     □ n>0 -> fac(i+1)!n-1;
//       r:integer; fac(i+1)?r; fac(i-1)!(n*r)
//   ]] || fac(0)::USER ]
//
// Note that the solution above from original paper is wrong.
// Check the code below for some fixes.
func S42_Factorial(fac []chan int, limit int) {
	for i := 1; i <= limit; i++ {
		go func(i int) {
			n := <-fac[i-1]
			if n == 0 {
				fac[i-1] <- 1
			} else if n > 0 {
				// Note that here we check if i equals limit.
				// The original solution in the paper fails to terminate
				// if user input is equal or higher than the given limit.
				if i == limit {
					fac[i-1] <- n
				} else {
					fac[i] <- n - 1
					r := <-fac[i]
					fac[i-1] <- n * r
				}
			}
		}(i)
	}
}

// S43_SmallSetOfIntegers implements Section 4.3 Small Set Of Integers.
// "To represent a set of not more than 100 integers as a process, S,
// which accepts two kinds of instruction from its calling process X:
// (1) S!insert(n), insert the integer n in the set and
// (2) S!has(n); ...; S?b, b is set true if n is in the set, and false
// otherwise. The initial value of the set is empty"
//
// Solution:
//
//   S::
//   content(0..99)integer; size:integer; size := 0;
//   *[n:integer; X?has(n) -> SEARCH; X!(i<size)
//   □ n:integer; X?insert(n) -> SEARCH;
//         [i<size -> skip
//         □i = size; size<100 ->
//            content(size) := n; size := size+1
//   ]]
//
// where SEARCH is an abbreviation for:
//
//   i:integer; i := 0;
//   *[i<size; conent(i) != n -> i:=i+1]
//
type S43_SmallSetOfIntegers struct {
	content []int
	size    int
}

// NewS43_SmallSetOfIntegers returns a S43_SmallSetOfIntegers
func NewS43_SmallSetOfIntegers() S43_SmallSetOfIntegers {
	return S43_SmallSetOfIntegers{content: make([]int, 100)}
}

// SEARCH returns the index of n if it is found in the set,
// otherwise returns the size of the set.
func (s *S43_SmallSetOfIntegers) SEARCH(n int) int {
	for i := 0; i < s.size; i++ {
		if s.content[i] != n {
			continue
		}
		return i
	}
	return s.size
}

// Has searches in the set given n, has receives true if n is found.
func (s *S43_SmallSetOfIntegers) Has(n int, has chan bool) {
	defer close(has)

	if s.SEARCH(n) < s.size {
		has <- true
		return
	}
	has <- false
	return
}

// Insert inserts given n, done recieves true if n is inserted.
func (s *S43_SmallSetOfIntegers) Insert(n int, done chan bool) {
	defer close(done)

	i := s.SEARCH(n)
	if i < s.size {
		done <- false
		return // nothing to do
	}
	// not found, insert to the array
	if i == s.size && s.size < 100 {
		s.content[s.size] = n
		s.size++
		done <- true
		return
	}

	done <- false
	return
}
