package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/shimin/pow/internal/server"
	"github.com/shimin/pow/internal/wisdom"
	"github.com/shimin/pow/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := LoadConfig()
	if err != nil {
		sugar.Fatal("cannot load config:", err)
	}

	wisdomSet, err := wisdom.NewSet("./WordsOfWisdom.json")
	if err != nil {
		sugar.Fatal("cannot load wisdom set:", err)
	}

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	ctx, done := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	h := server.NewHandler(sugar, cfg.KeySize, cfg.TargetBits, wisdomSet)
	serverOptions := []grpc.ServerOption{grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    cfg.GrpcPingTime,
		Timeout: cfg.GrpcTimeOut,
	})}
	grpcServer := grpc.NewServer(serverOptions...)
	proto.RegisterAuthServiceServer(grpcServer, h)

	l, err := net.Listen("tcp", cfg.Host)
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

	sugar.Infof("Server is listening at %s", cfg.Host)
	wg.Wait()
}
