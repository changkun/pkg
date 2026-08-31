# pkg

[![CI](https://github.com/changkun/pkg/actions/workflows/ci.yml/badge.svg)](https://github.com/changkun/pkg/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/changkun.de/x/pkg.svg)](https://pkg.go.dev/changkun.de/x/pkg)
[![Go Report Card](https://goreportcard.com/badge/changkun.de/x/pkg)](https://goreportcard.com/report/changkun.de/x/pkg)

A personal Go codebase: data structures, concurrency utilities, numerical
code, and small experiments that came out of other work. Nothing here promises
a stable API. Use it at your own risk.

```
go get changkun.de/x/pkg
```

Go 1.27 or later.

## Packages

### Data structures and algorithms

| Package | What it gives you |
| --- | --- |
| [`ds`](ds) | LRU cache, red-black tree, skip list, ring buffer, queue, stack, min stack, randomized set |
| [`algo`](algo) | Binary search, merge sort, insertion sort over a caller-supplied comparator |
| [`slice`](slice) | Slice helpers |
| [`mat`](mat) | Dense matrices with cache-aware and parallel multiplication kernels, and eigen decomposition |

### Concurrency

| Package | What it gives you |
| --- | --- |
| [`lockfree`](lockfree) | Lock-free queue and stack, and atomic float64 arithmetic |
| [`promise`](promise) | Wait group with a timeout, and a retry strategy |
| [`sched`](sched) | Task scheduler you can pause, resume and stop |
| [`g`](g) | Goroutine lifecycle management, with cancellation and undo |
| [`gls`](gls) | Goroutine local storage |
| [`csp`](csp) | The examples of Hoare's 1978 CSP paper, written in Go |
| [`maint`](maint) | Run a function on the main thread |
| [`leaktest`](leaktest) | Fail a test that leaves goroutines behind |

### Runtime and system

| Package | What it gives you |
| --- | --- |
| [`rt`](rt) | Goroutine id, function names, and a GC notification signal |
| [`mkill`](mkill) | Cap the number of OS threads the runtime keeps alive |
| [`sysmon`](sysmon) | Watch a load signal and scale a resource before it runs out |
| [`gen`](gen) | Fast random numbers and random strings |

### Numerical and signal

| Package | What it gives you |
| --- | --- |
| [`bo`](bo) | Bayesian optimization over a Gaussian process |
| [`rp`](rp) | Peak load detection: a z test for a rising trend, and a Poisson acceptance test |
| [`detect`](detect) | Moving average and Kolmogorov-Zurbenko filters |
| [`convert`](convert) | Linear to sRGB colour conversion |

### Services and I/O

| Package | What it gives you |
| --- | --- |
| [`net`](net) | One-call JSON HTTP requests with a context and a timeout |
| [`mailgun`](mailgun) | Send an email through Mailgun |
| [`hue`](hue) | Control Philips Hue lights, as a library and an `office` command |
| [`caption`](caption) | Speech to text through Google Cloud Speech and xfyun |
| [`errors`](errors) | Try and catch style error handling |

### Experiments

[`benchs`](benchs) compares REST against gRPC, cgo against Go syscalls, and
call overheads. [`ray`](ray) is a ray tracer. [`metal`](metal), [`win`](win)
and [`hotkey`](hotkey) are graphics and input experiments and only build on
macOS. [`mem`](mem) measures huge page prefetching and only builds on Linux.
[`misc`](misc) holds standalone examples.

## Development

```
go build ./...                # builds everything for the current platform
go test -race ./...           # tests
golangci-lint run ./...       # lint, pinned to the version CI uses
```

CI runs the tests and the linter on both Linux and macOS, because several
packages only build on one of them, and runs `govulncheck` on every push and
weekly. The protobuf definitions under `benchs/restrpc/rpcs` are regenerated
with `go generate ./benchs/restrpc/rpcs/`, which uses the tools pinned in
`go.tool.mod` and needs no system `protoc`.

## License

MIT &copy; [changkun](https://changkun.de)
