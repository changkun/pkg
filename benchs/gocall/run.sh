#!/bin/bash
REPEAT=10

gcc -O2 c/ccall.c -o ccall
go test -bench=. -count=$REPEAT

for i in `seq 1 $REPEAT`; do
    ./ccall
done
rm ccall