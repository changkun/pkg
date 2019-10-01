// Package mat is a golang matrix package with cache-aware lock-free
// tiling optimization.
//
// The following contents explains how multiplication algorithm works.
//
// Prior knowledge
//
// Assume only 2 levels in the hierarchy, fast(registers/cache) and slow(main memory).
// All data initially in slow memory
//
//  m:       number of memory elements (words) moved between fast and slow memory
//  t_m:     time per slow memory operation
//  f:       number of arithemetic operations
//  t_f:     time per arithmetic operation (t_f << t_m)
//  q = f/m: computational intensity (key to algorithm efficiency) average number
//           of flops per slow memory access
//
//  Minimum possible time = f * t_f when all data in fast memory.
//  Actual time = f * t_f + m * t_m = f * t_f * [1 + (t_m / t_f) * (1 / q)]
//  Machine balance a = t_m / t_f (key to machine efficiency)
//
// Larger q means time closer to minimum f*t_f
//
//  q >= t_m / t_f needed to get at least half of peak speed
//
// Tiled matrix multiply
//
//  C = A·B =
//   /          \   /         \      /                                \
//  |  A11  A12  | |  B11 B12  |    | A11·B11+A12·B21  A11·B12+A12·B22 |
//  |            | |           | == |                                  |
//  |  A21  A22  | |  B21 B22  |    | A21·B11+A22·B22  A21·B12+A22·B22 |
//   \          /   \         /      \                                /
//
// Consider A, B, C to be N-by-N matrices of b-by-b sub-blocks:
//
//  b = n / N is called the *block size*
//
// Thus
//
//  m = N*n*n read each block of B N*N*N times (N*N*N * b*b = N*N*N * (n/N)*(n/N) = N*n*n)
//    + N*n*n read each block of A N*N*N times
//    + 2*n*n read and write each block of C once
//    = 2*(N+1)*n*n
//
// Hence, computational intensity q = f/m = 2*n*n*n / (2*(N+1)*n*n)
//
//  q ~= n/N = b (for large n)
//
// So we can improve performance by increasing the blocksize b
// Can be much faster than matrix-vector multiply (q=2)
//
// The blocked algorithm has computational intensity q ~= b: the large
// the block size, the more efficient algorithm will be
//
// Limit: All three blocks from A, B, C must fit in fast memory (cache),
// so we cannot make these blocks arbitrarily large.
//
// Assume our fast memory has size M_fast:
//
//  3*b*b <= M_fast, so q ~= b <= sqrt(M_fast / 3)
//
// To build a machine to run matrix multiply at 1/2 peak arthmetic speed
// of the machine, we need a fast memory of size
//
//  M_fast >= 3*b*b ~= 3*q*q = 3*(t_m / t_f)*(t_m / t_f)
//
// This size is reasonable for L1 cache, but not for register sets
// Note: analysis assumes it is possible to schedule the instructions
// perfectly.
//
// Limits
//
// The tiled algorithm changes the order in which values are accumulated
// into each C[i, j] by applying commutativity and associativity: Get slightly
// different answers from naive version, because of roundoff.
//
// The previous anlysis showed that the blocked algorithm has computational
// intensity:
//
//  q ~= b <= sqrt(M_fast / 3)
//
// There is a lower bound result that says we cannot do any better than this
// (using only associativity):
//
// Theorem (Hong & Kong, 1981): Any reorganization of this algorithm (that uses
// only associativity) is limited to:
//
//  q = O(sqrt(M_fast))
package mat
