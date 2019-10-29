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
// Monitors and Scheduling
//
// A monitor is prepared to communicate with any of its user processes (
// i.e. whichever of them calls first) it will use a guarded command
// with a range.
//
// Author: Changkun Ou <hi@changkun.us>
package csp

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

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

// S44_Scan implements Section 4.4 Scanning a Set
// "Extend the solution to 4.3 by providing a fast method for scanning
// all members of the set without changing the value of the set. The
// user program will contain a repetitive command of the form."
//
// Solution:
//
//   S!scan(); more:boolean; more:=true;
//   *[more;x:integer;S?next(x)->...deal with x ... .
//   □ more;S?noneleft()->more:=false]
//
func (s *S43_SmallSetOfIntegers) S44_Scan(recv chan int) {
	for _, v := range s.content {
		recv <- v
	}
	close(recv)
}

// S45_RecursiveSmallSetOfIntegers implements Section 4.5 Recursive
// Data Representation: Small Set of Integers.
// "Same as above, but an array of processes is to be used to achieve a
// high degree of parallelism. Each process should contain at most one
// number. When it contains no number, it should answer 'false' to all
// inquiries about membership. On the first insertion, it changes to a
// second phase of behavior, in which it deals with instructions from
// its predecessor, passing some of them on to its successor. The
// calling process will be named S(0). For efficiency, the set should be
// sorted, i.e. the i-th process should contain the i-th largest number."
//
// Solution:
//
//   S(i:1..100)::
//   *[n:integer; S(i-1)?has(n)->S(0)!false
//   □ n:integer; S(i-1)?insert(n)->
//      *[m:integer; S(i-1)?has(m)->
//         [m<=n->S(0)!(m=n)
//         □m>n-->S(i+1)!has(m)
//       ]
//      □m:integer; S(i-1)?insert(m)->
//       [m<n->S(i-1)!insert(m); n:=m
//       □m=n->skip
//       □m>n->S(i+1)!insert(m)
//   ]]]
//
// S46_RemoveTheLeastMember implements Section 4.6 Multiple Exits:
// Remove the Least Member.
// "Exercise: Extend the above solution to respond to a command to yield
// the least member of the set and to remove it from the set. The user
// program will invoke the facility by a pair of commands:
//
//  S(1)!least();[x:integer;S(1)?x-> ... deal with x ...
//               □S(1)?nonleft()-> ... ]
//
// or if he wishes to scan and empty the set, he may write:
//
//  S(1)!least();more:boolean;more:=true;
//               *[more;xinteger;S(1)?x-> ...deal with x...; S(1)!least())
//                □ more;S(1)?noneleft()->more:=false]
// "
func S45_S46_NewRecursiveSmallSetOfIntegers() (chan S45_Has, chan int, chan chan S46_Least) {
	// FIXME: this implementation doesn't close all processes
	size := 100
	has := make([]chan S45_Has, size)
	insert := make([]chan int, size)
	least := make([]chan S46_Least, size)
	leastQuery := make([]chan chan S46_Least, size)
	for i := 0; i < size; i++ {
		has[i] = make(chan S45_Has)
		insert[i] = make(chan int)
		least[i] = make(chan S46_Least)
		if i == 0 {
			leastQuery[i] = make(chan chan S46_Least)
		}

		go func(i int) {
			// a goroutine that stores the actual value
			var n int

		EMPTY:
			for {
				select {
				case h := <-has[i]:
					h.Response <- false
				case n = <-insert[i]:
					goto NONEMPTY
				case least[i] <- S46_Least{NoneLeft: true}:
					continue
				case q := <-leastQuery[i]:
					q <- S46_Least{NoneLeft: true}
				}
			}
		NONEMPTY:
			for {
				select {
				case h := <-has[i]:
					if h.V <= n {
						h.Response <- h.V == n
					} else {
						if i == size {
							h.Response <- false
						} else {
							has[i+1] <- h
						}
					}
				case in := <-insert[i]:
					if in < n {
						insert[i+1] <- in
						n = in
					} else if in == n {
						continue
					} else if in > n {
						insert[i+1] <- in
					}
				case least[i] <- S46_Least{n, false}:
					next := <-least[i+1]
					if next.NoneLeft {
						goto EMPTY
					} else {
						n = next.Least
					}
				case l := <-leastQuery[i]:

					next := <-least[i+1]
					l <- S46_Least{n, next.NoneLeft}
					if next.NoneLeft {
						goto EMPTY
					} else {
						n = next.Least
					}
				}
			}
		}(i)
	}
	return has[0], insert[0], leastQuery[0]
}

type S45_Has struct {
	V        int
	Response chan bool
}

type S46_Least struct {
	Least    int
	NoneLeft bool
}

// S51_BoundedBuffer implements Section 5.1 Bounded Buffer
// "Construct a buffering process X to smooth variations in the speed of
// output of portions by a producer process and input by a consumer
// process. The consumer contains pairs of commands X!more();X?p, and
// the producer contains commands of the form X!p. The buffer should
// contain up to ten portions."
//
// Solution:
//
//   X::
//   buffer:(0..9)portion;
//   in,out:integer; in:=0; out := 0;
//   comment 0 <= out <= in <= out+10;
//     *[in < out + 10; producer?buffer(in mod 10) -> in := in + 1
//     □ out < in; consumer?more() -> consumer!buffer(out mod 10);
//        out := out + 1 ]
func S51_BoundedBuffer() (chan int, chan int) {
	in, out := 0, 0
	size := 10
	buffer := make([]int, size)
	producer, consumer := make(chan int), make(chan int)

	go func() {
		for {
			if in < out+size {
				select {
				case v := <-producer:
					buffer[in%size] = v
					in++
				default:
				}
			}
			if out < in {
				select {
				case consumer <- buffer[out%size]:
					out++
				default:
				}
			}
			runtime.Gosched()
		}
	}()
	return producer, consumer
}

// S52_IntegerSemaphore implements Section 5.2 Integer Semaphore.
// "To implement an integer semaphore, S, shared among an array
// X(i:1..100) of client processes. Each process many increment the
// semaphore by S!V() or delayed if the value of the semaphore is not
// positive."
//
// Solution:
//
//   S::val:integer; val:=0;
//   *[(i:1..100)X(i)?V()->val:=val+1
//   □ (i:1..100)val>0;X(i)?P()->val:=val-1]
type S52_IntegerSemaphore struct {
	inc  chan struct{}
	dec  chan struct{}
	done chan struct{}
}

func NewS52_IntegerSemaphore() S52_IntegerSemaphore {
	sem := S52_IntegerSemaphore{
		inc:  make(chan struct{}),
		dec:  make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		var n int
		select {
		case <-sem.inc:
			n++
		case <-sem.done:
			return
		}
		for {
			select {
			case <-sem.inc:
				n++
			case <-sem.dec:
				n--
				// block until next increase
				if n == 0 {
					select {
					case <-sem.inc:
						n++
					case <-sem.done:
						return
					}
				}
			case <-sem.done:
				return
			}
		}
	}()
	return sem
}

// P operator decreases semaphore
func (s *S52_IntegerSemaphore) P() {
	s.dec <- struct{}{}
}

// V operator increases semaphore
func (s *S52_IntegerSemaphore) V() {
	s.inc <- struct{}{}
}

// Close closes the semaphore
func (s *S52_IntegerSemaphore) Close() {
	close(s.done)
}

// S53_DiningPhilosophers implements Section 5.3 Dining Philosophers
// "Five philosophers spend their lives thinking and eating. The
// philosophers share a common dining room where there is a curcular
// table surrounded by five chairs, each belonging to one philosopher.
// In the center of the table there is a large bowl of spaghetti, and
// the table is laid with five forks. On feeling hungry, a philosopher
// enters the dinning room, sits in his own chair, and picks up the fork
// on the left of his place. Unfortunately, the spaghetti is so tangled
// that he needs to pick up and use the fork on his right as well. When
// he has finished, the puts down both forks, and leaves the room. The
// room should keep a count of the number of philosophers in it."
//
// Solution:
//
// The behavior of the i-th philosopher may be described as follows:
//
//   PHIL = *[...during ith lifetime ... ->
//            THINK;
//            room!enter();
//            fork(i)!pickup();fork((i+1)mod5)!pickup();
//            EAT;
//            fork(i)!putdown();fork((i+1)mod5)!putdown();
//            room!next()]
//
// The fate of i-th fork is to be picked up and put down by a
// philosopher sitting on either side of it
//
//   FORK = *[phil(i)?pickup()->phil(i)?putdown()
//          □ phil((i-1)mod5)?pickup()->phil((i-1)mod5)?putdown()]
//
// The story of the room may be simply told:
//
//   ROOM = occupancy:integer; occupancy := 0;
//          *[(i:0..4)phil(i)?enter()->occupancy:=occupancy+1
//          □ (i:0..4)phil(i)?exit()->occupancy:=occupancy-1]
//
// All these components operate in parallel:
//
//   [room::ROOM||fork(i:0..4)::FORK||phil(i:0..4)::PHIL]
func S53_DiningPhilosophers() {
	// FIXME: may deadlock? mentioned in the exercise
	size := 5
	enter := make(chan int)
	exit := make(chan int)
	pickup := make([]chan struct{}, size)
	putdown := make([]chan struct{}, size)
	for i := 0; i < size; i++ {
		pickup[i] = make(chan struct{})
		putdown[i] = make(chan struct{})
	}
	THINK := func(i int) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
	}
	EAT := func(i int) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
	}
	PHIL := func(i int) {
		THINK(i)
		enter <- i
		pickup[i] <- struct{}{}
		pickup[(i+1)%5] <- struct{}{}
		EAT(i)
		putdown[i] <- struct{}{}
		putdown[(i+1)%5] <- struct{}{}
		exit <- i
	}
	FORK := func(i int) {
		for {
			select {
			case <-pickup[i]:
				<-putdown[i]
			case <-pickup[(i+4)%5]:
				<-putdown[(i+4)%5]
			}
		}
	}
	ROOM := func() {
		occupancy := 0
		for {
			select {
			case i := <-enter:
				occupancy++
				fmt.Printf("%v enter room, occupancy: %v\n", i, occupancy)
			case i := <-exit:
				occupancy--
				fmt.Printf("%v exit room, occupancy: %v\n", i, occupancy)
			}
		}
	}
	go ROOM()
	wg := sync.WaitGroup{}
	wg.Add(size)
	for i := 0; i < size; i++ {
		go FORK(i)
		go func(i int) {
			PHIL(i)
			wg.Done()
		}(i)
	}
	wg.Wait() // wait until all philosophers are finished.
}

// S61_TheSieveOfEratosthenes implements Section 6.1 Prime Numbers: The
// Sieve of Eratosthenes.
// "To print in ascending order all primes less than 10000. Use an array
// of processes, SIEVE, in which each process inputs a prime from its
// predecessor and prints it. The process then inputs an ascending
// stream of numbers from its predecessor and passes them on to its
// successor, suppressing any that are multiples of the original prime."
//
// Solution:
//
//   [SIEVE(i:1..100)::
//    p,mp:integer;
//    SIEVE(i-1)?p;
//    print!p;
//    mp:=p; comment mp is a multiple of p;
//   *[m:integer; SIEVE(i-1)?m->
//      *[m>mp->mp:=mp+p];
//       [m=mp->skip □ m<mp->SIEVE(i+1)!m ]
//    ]
//   || SIEVE(0)::print!2; n:integer; n:=3;
//         *[n<10000->SIEVE(1)!n;n:=n+2]
//   || SIEVE(101)::*[n:integer;SIEVE(100)?n->print!n]
//   || print::*[(i:0..101)n:integer;SIEVE(i)?n->...]
//   ]
func S61_TheSieveOfEratosthenes(np int) {

	SIEVE := make([]chan int, np+1)
	for i := 1; i <= np; i++ {
		SIEVE[i] = make(chan int)

		go func(i int) {
			p, ok := <-SIEVE[i]
			if !ok {
				return
			}
			println(p)

			mp := p // mp is a multiple of p
			for {
				m := <-SIEVE[i]
				for m > mp {
					mp += p
				}
				if m < mp {
					if i < np {
						SIEVE[i+1] <- m
					} else {
						// this fixes a bug in original paper
						SIEVE[1] <- m
					}
				}
			}
		}(i)
	}
	done := make(chan bool)

	go func() {
		println(2)
		for n := 3; n < 10000; n += 2 {
			SIEVE[1] <- n
		}
		// FIXME: send n finished does not meaning that all primes are
		// printed, because of the original algorithm in the paper is
		// incorrect, our quick fix doesn't guarantee the order of
		// outputs in ascending order.
		done <- true
	}()

	for {
		select {
		case p, ok := <-SIEVE[100]:
			if !ok {
				return
			}
			println(p)
		case <-done:
			return
		}
	}
}

// S62_MatrixMultiplication implements Section 6.2 An interative Array:
// Matrix Multiplication.
// "A square matrix A of order 3 is given. Three streams are to be input,
// each stream representing a column of an array IN. There streams are
// to be output, each representing a column of the product matrix IN x A.
// After an initial delay, the results are to be produced at the same
// rate as the input is consumed. Consequently, a high degree of
// parallelism is required. Each of the nine nonborder nodes inputs a
// vector components from the west and a partial sum from the north.
// Each node outputs the vector component to its east, and an updated
// partial sum to the south. The input data is produced by the west
// border nodes, and the desired results are consumed by south border
// nodes. The north border is a constant source of zeros and the east
// border is just a sink. No provision need be made for termination nor
// for changing the values of the array A."
//
// Solution: There are twenty-one nodes, in five groups, comparising the
// central square and the four borders:
//
//   [M(i:1..3,0)::WEST
//   ||M(0,j:1..3)::NORTH
//   ||M(i:1..3,4)::EAST
//   ||M(4,j:1..3)::SOUTH
//   ||M(i:1..3,j:1..3)::CENTER]
//
// The WEST and SOUTH borders are processes of the user program; the
// remaining processes are:
//
//   NORTH = *[true -> M(1,j)!0]
//   EAST = *[x:real; M(i,3)?x->skip]
//   CENTER = *[x:real;M(i,j-1)?x->
//             M(i,j+1)!x;sum:real;
//             M(i-1,j)?sum;M(i+1,j)!(A(i,j)*x+sum)]
type S62_Matrix struct {
	WEST  []chan int // for input
	SOUTH []chan int // for output

	in  [][]chan int
	inA [][]chan int

	A [][]int
}

func (m *S62_Matrix) NORTH(j int) {
	for {
		m.inA[0][j] <- 0.0
	}
}

func (m *S62_Matrix) EAST(i int) {
	for {
		println("skip:", <-m.in[i][len(m.in[i])-1]) // skip command
	}
}

func (m *S62_Matrix) CENTER(i, j int) {
	for x := range m.in[i][j-1] {
		m.in[i][j] <- x
		sum := <-m.inA[i-1][j]
		fmt.Printf("m.inA[%d][%d]: %v\n", i-1, j-1, m.A[i-1][j-1]*x+sum)
		m.inA[i][j] <- m.A[i-1][j-1]*x + sum
	}
}

func S62_NewMatrix(A [][]int) S62_Matrix {
	m := S62_Matrix{
		WEST: make([]chan int, 3),
		in:   make([][]chan int, 3+1),
		inA:  make([][]chan int, 3+1),
		A:    A,
	}

	for i := 0; i < 3+1; i++ {
		m.in[i] = make([]chan int, 3+1)
		m.inA[i] = make([]chan int, 3+1)
		for j := 0; j < 3+1; j++ {
			m.in[i][j] = make(chan int)
			m.inA[i][j] = make(chan int)
		}
	}

	for i := 1; i <= 3; i++ {
		m.WEST[i-1] = m.in[i][0]
	}
	m.SOUTH = m.inA[3][1:]

	for i := 1; i <= 3; i++ {
		go m.NORTH(i)
		go m.EAST(i)
		for j := 1; j <= 3; j++ {
			go m.CENTER(i, j)
		}
	}

	return m
}

func (m *S62_Matrix) S62_Multiply(IN [][]int) (OUT [][]int) {
	for i := 0; i < 3; i++ {
		go func(i int) {
			for j := 0; j < 3; j++ {
				m.WEST[j] <- IN[i][j]
			}
		}(i)
	}

	OUT = make([][]int, 3)
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		OUT[i] = make([]int, 3)
		go func(i int) {
			for j := 0; j < 3; j++ {
				OUT[i][j] = <-m.SOUTH[i]
			}
			wg.Done()
		}(i)
	}
	println(1)
	wg.Wait()
	return
}
