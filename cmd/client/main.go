package main

import (
	"context"
	"encoding/binary"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/shimin/pow/internal/pow"
	"github.com/shimin/pow/internal/utils"
	"github.com/shimin/pow/proto"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := LoadConfig()
	if err != nil {
		sugar.Fatal("cannot load config:", err)
	}

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	ctx, done := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case <-sigquit:
			done()
		case <-ctx.Done():
			return
		}
	}()

	connector := utils.NewGrpcClientConnector(sugar, cfg.Host, cfg.Retries, cfg.RetryInterval, cfg.ConnTimeout, cfg.ConnTimeout)
	conn, err := connector.GrpcConnectWithRetry(ctx)

	if err != nil && ctx.Err() == nil {
		sugar.Errorf("unable to establish connection")
		return
	}

	if ctx.Err() != nil {
		return
	}

	defer conn.Close()

	c := proto.NewAuthServiceClient(conn)

	go func(client proto.AuthServiceClient) {
		defer wg.Done()
		defer done()

		stream, err := client.AuthFlow(ctx)
		if err != nil {
			sugar.Error(err)
			return
		}

		msg, err := stream.Recv()
		if err != nil {
			sugar.Error(err)
			return
		}

		data := msg.GetData()

		if len(data) == 0 {
			sugar.Error("Empty data")
			return
		}

		result := pow.Calc(ctx, data, cfg.TargetBits)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, result)

		stream.Send(&proto.Packet{
			Data: buf,
		})

		if err != nil {
			sugar.Error(err)
			return
		}

		msg, err = stream.Recv()
		if err != nil {
			sugar.Error(err)
			return
		}

		sugar.Info(string(msg.GetData()))
	}(c)

	wg.Wait()
}
