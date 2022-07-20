package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ishimin/antidos/internal/pow"
	"go.uber.org/zap"
)

const (
	host       = ":5000"
	keySize    = 40
	targetBits = 24
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", host)
	if err != nil {
		sugar.Fatal(err)
	}

	ctx, done := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for ctx.Err() == nil {
			conn, err := listener.Accept()
			if err != nil {
				sugar.Warn(err)
				continue
			}
			go handleConnection(ctx, conn, sugar)
		}
	}()

	go func() {
		defer wg.Done()
		<-sigquit
		done()
		listener.Close()
	}()

	sugar.Infof("Server is listening at %s", host)
	wg.Wait()
}

func handleConnection(ctx context.Context, conn net.Conn, log *zap.SugaredLogger) {
	defer conn.Close()
	data := make([]byte, keySize)
	rand.Read(data)
	_, err := conn.Write(data)
	if err != nil {
		log.Error(err)
		return
	}

	ans := make([]byte, 8)
	n, err := conn.Read(ans)
	if err != nil {
		log.Error(err)
		return
	}
	if n != 8 {
		log.Error("size answer mismatch")
		return
	}

	ok := pow.Validate(data, targetBits, binary.LittleEndian.Uint64(ans))
	if ok {
		fmt.Fprint(conn, "authorised")
		return
	}

	fmt.Fprint(conn, "access denied")
}
