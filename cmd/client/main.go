package main

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ishimin/antidos/internal/pow"
	"github.com/ishimin/antidos/internal/utils"
)

// const targetBits = 24

// func main() {
// 	result := pow.Calc(context.Background(), "some string", targetBits)
// 	fmt.Printf("Pow: %t\n", pow.Validate("some string", targetBits, result))
// 	fmt.Println()
// }

// move to config
const (
	host          = "127.0.0.1:5000"
	targetBits    = 24
	retries       = 2
	retryInterval = 2 * time.Second
	connTimeout   = 2 * time.Second
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

	go func() {
		defer wg.Done()

		client := utils.NewClientConnector(sugar, host, retries, retryInterval, connTimeout)
		conn, err := client.ConnectWithRetry(ctx)
		if err != nil {
			sugar.With(err).Errorf("unable to establish connection")
			done()
			return
		}

		data := make([]byte, 40)
		_, err = conn.Read(data)

		if err != nil {
			sugar.Error(err)
			conn.Close()
			return
		}

		if len(data) == 0 {
			sugar.Error("Empty request")
			conn.Close()
			return
		}

		result := pow.Calc(ctx, data, targetBits)
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, result)

		_, err = conn.Write(buf)
		if err != nil {
			sugar.Error(err)
			return
		}

		d, err := ioutil.ReadAll(conn)
		if err != nil {
			sugar.Error(err)
			return
		}

		sugar.Info(string(d))
		done()
	}()

	go func() {
		defer wg.Done()
		select {
		case <-sigquit:
			done()
		case <-ctx.Done():
			return
		}
	}()

	wg.Wait()
}
