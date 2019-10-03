// Package csp implements "Hoare, C. A. R. (1978). Communicating
// sequential processes. Communications of the ACM, 21(8), 666–677.
// https://doi.org/10.1145/359576.359585".
//
// hi <at> changkun.us
package csp

// S31_COPY implements Section 3.1 COPY problem:
// "Write a process X to copy characters output by process west to
// process, east."
func S31_COPY(west, east chan rune) {
	for c := range west {
		east <- c
	}
	close(east)
}
