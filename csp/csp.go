// Package csp implements "Hoare, C. A. R. (1978). Communicating
// sequential processes. Communications of the ACM, 21(8), 666–677.
// https://doi.org/10.1145/359576.359585".
//
// hi <at> changkun.us
package csp

// S31_COPY implements Section 3.1 COPY problem:
// "Write a process X to copy characters output by process west to
// process, east."
//
// Solution:
// X :: *[c:character; west>c -> east!c]
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
// X :: *[c:character; west?c ->
//   [ c != asterisk -> east!c
//    □ c = asterisk -> west?c;
//          [ c != asterisk -> east!asterisk; east!c
//           □ c = asterisk -> east!upward arrow
//   ] ]    ]
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

// S33_DISASSEMBLE implements Section 3.3 DISASSEMBLE problem:
// "to read cards from a cardfile and output to process X the stream of
// characters they contain. An extra space should be inserted at the end
// of each card."
//
// Solution:
// *[cardimage:(1..80)characters; cardfile?cardimage ->
//     i:integer; i := 1;
//     *[i <= 80 -> X!cardimage(i); i := i+1 ]
//     X!space
// ]
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
