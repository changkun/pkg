// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ser

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"changkun.de/x/pkg/benchs/restrpc/route"
	"changkun.de/x/pkg/benchs/restrpc/rpcs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	addrHTTP        = "0.0.0.0:12345"
	addrGRPC        = "0.0.0.0:12346"
	maxMessageSize  = 500 << 20 // 500 MB
	shutdownTimeout = 5 * time.Second
	connTimeout     = 5 * time.Minute
)

// RunHTTP serves the REST endpoint until the process is interrupted. It
// returns the first error that stopped the server, and nil on a clean
// shutdown.
func RunHTTP() error {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	server := &http.Server{
		Handler:           route.Register(),
		Addr:              addrHTTP,
		ReadHeaderTimeout: connTimeout,
	}

	shutdown := make(chan error, 1)
	go func() {
		// os.Kill cannot be caught, so listening for it never fires.
		// SIGTERM is the signal a supervisor actually sends.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		shutdown <- server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return <-shutdown
}

// RunRPC serves the gRPC endpoint until Serve returns.
func RunRPC() error {
	l, err := net.Listen("tcp", addrGRPC)
	if err != nil {
		return err
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
		grpc.ConnectionTimeout(connTimeout),
	)
	rpcs.RegisterArithmeticServer(s, &rpcs.Server{})
	return s.Serve(l)
}
