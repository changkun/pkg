// Package csp implements "Hoare, C. A. R. (1978). Communicating
// sequential processes. Communications of the ACM, 21(8), 666–677.
// https://doi.org/10.1145/359576.359585".
//
// Author: Changkun Ou <hi@changkun.us>
package csp

// S31_COPY implements Section 3.1 COPY problem:
// "Write a process X to copy characters output by process west to
// process, east."
//
// Solution:
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
	cache := make([]rune, 125)

	i := 0
	for c := range X {
		cache[i] = c
		if i <= 124 {
			i++
		}
		if i == 125 {
			cache[i-1] = c
			lineprinter <- string(cache)
			i = 0
		}
	}
	if i > 0 {
		for i <= 124 {
			cache[i] = ' '
			i++
		}
		lineprinter <- string(cache)
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
//   [west::DISASSEMBLE||X::SQUASH||east::ASSEMBLE]
func S36_ConwayProblem(cardfile chan []rune, lineprinter chan string) {
	west, east := make(chan rune), make(chan rune)
	go S33_DISASSEMBLE(cardfile, west)
	go S32_SQUASH_EX(west, east)
	S34_ASSEMBLE(east, lineprinter)
}
