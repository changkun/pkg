#!/bin/bash
REPEAT=5
gcc -O2 c/syscall.c -o syscall_c
go test -bench=. -count=$REPEAT -timeout 20m -v
for i in `seq 1 ${REPEAT}`; do
    ./syscall_c
done
rm syscall_c