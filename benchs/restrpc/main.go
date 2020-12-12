//go:generate protoc -I rpcs rpcs/add.proto --go_out=plugins=grpc:rpcs

package main

import "changkun.de/x/pkg/benchs/restrpc/ser"

func main() {
	go ser.RunRPC()
	ser.RunHTTP()
}
