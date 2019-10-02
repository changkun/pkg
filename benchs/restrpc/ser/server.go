package ser

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/changkun/gobase/benchs/restrpc/route"
	"github.com/changkun/gobase/benchs/restrpc/rpcs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	addrHTTP       = "0.0.0.0:12345"
	addrGRPC       = "0.0.0.0:12346"
	maxMessageSize = 500 << 20 // 500 MB
)

// RunHTTP ...
func RunHTTP() {
	gin.DefaultWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)
	server := &http.Server{
		Handler: route.Register(),
		Addr:    addrHTTP,
	}
	terminated := make(chan bool, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := server.Shutdown(ctx); err != nil {
			panic(err)
		}

		cancel()
		terminated <- true
	}()

	// logrus.Infof("http is running")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
	<-terminated
}

// RunRPC ...
func RunRPC() {
	l, err := net.Listen("tcp", addrGRPC)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.MaxMsgSize(maxMessageSize),
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxSendMsgSize(maxMessageSize),
		grpc.ConnectionTimeout(time.Minute*5),
	)
	rpcs.RegisterArithmeticServer(s, &rpcs.Server{})
	// logrus.Infof("grpc is running")
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
