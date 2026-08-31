// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

//go:generate go tool -modfile=../../../go.tool.mod buf generate

package rpcs

import "context"

// Server answers the Arithmetic service. It embeds the generated
// UnimplementedArithmeticServer, which grpc-go requires so that a service
// gaining a method does not break every implementation of it.
type Server struct {
	UnimplementedArithmeticServer
}

// Add returns the sum of the two operands.
func (s *Server) Add(ctx context.Context, in *AddInput) (*AddOutput, error) {
	return &AddOutput{Sum: in.A + in.B}, nil
}
