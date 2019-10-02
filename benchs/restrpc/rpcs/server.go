package rpcs

import "context"

// Server ...
type Server struct{}

// Add ...
func (s *Server) Add(ctx context.Context, in *AddInput) (out *AddOutput, err error) {
	return &AddOutput{Sum: in.A + in.B}, nil
}
