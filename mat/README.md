# mat

Matrix package with cache-aware lock-free tiling optimization.

## Getting started

The following illustrates some basic usage of `mat`.

```go
// Create 2x3 matrix, specified its value
// New() will throw error if provided values 
// is ineuqal to its dimension
A, err := mat.New(2, 3)(
    1, 2, 3,
    4, 5, 6,
)

// Create a 3x4 random matrix
B := mat.Rand(3, 4)

// Create a 2x4 zero matrix
C := mat.Zero(2, 4)


// C = A x B, throw err if dimentionality error
err = A.Dot(B, C)
err = A.Dot(B, C) // with concurrency optimization
// or
C, err := mat.Dot(A, B) // alloc new matrix inside Dot
```

## License

MIT &copy; [changkun](https://changkun.de)