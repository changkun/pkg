// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Command restrpc serves the same addition over REST and over gRPC so the two
// transports can be benchmarked against each other.
package main

import (
	"log"

	"changkun.de/x/pkg/benchs/restrpc/ser"
)

func main() {
	go func() {
		if err := ser.RunRPC(); err != nil {
			log.Fatalf("grpc server: %v", err)
		}
	}()
	if err := ser.RunHTTP(); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
