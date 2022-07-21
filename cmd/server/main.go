package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shimin/pow/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	host         = ":5000"
	keySize      = 40
	targetBits   = 24
	grpcPingTime = 30 * time.Second
	grpcTimeOut  = 60 * time.Second
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	ctx, done := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	h := NewHandler(sugar)
	serverOptions := []grpc.ServerOption{grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    grpcPingTime,
		Timeout: grpcTimeOut,
	})}
	grpcServer := grpc.NewServer(serverOptions...)
	proto.RegisterAuthServiceServer(grpcServer, h)

	l, err := net.Listen("tcp", host)
	if err != nil {
		sugar.With(err).Fatalf("server can't listen and serve requests")
	}

	go func(ctx context.Context, srv *grpc.Server, listener net.Listener) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			l.Close()
			sugar.Infof("grpc listener closed")
			return
		default:
			if err = srv.Serve(listener); err != nil {
				sugar.With(err).Fatalf("server can't listen and serve requests")
			}
		}
	}(ctx, grpcServer, l)

	go func() {
		defer wg.Done()
		<-sigquit
		done()
		grpcServer.Stop()
		sugar.Infof("grpc listener closed")
	}()

	sugar.Infof("Server is listening at %s", host)
	wg.Wait()
}
