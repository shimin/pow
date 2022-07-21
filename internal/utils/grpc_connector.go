package utils

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GrpcClientConnector struct {
	log           *zap.SugaredLogger
	host          string
	retries       int
	timeout       time.Duration
	retryInterval time.Duration
	grpcPingTime  time.Duration
}

func NewGrpcClientConnector(log *zap.SugaredLogger,
	host string,
	retries int,
	timeout,
	retryInterval,
	grpcPingTime time.Duration) GrpcClientConnector {
	return GrpcClientConnector{
		log:           log,
		host:          host,
		retries:       retries,
		timeout:       timeout,
		retryInterval: retryInterval,
		grpcPingTime:  grpcPingTime,
	}
}

func (c *GrpcClientConnector) GrpcConnectWithRetry(ctx context.Context) (*grpc.ClientConn, error) {
	for i := 1; i <= c.retries && ctx.Err() == nil; i++ {
		c.log.Infof("connecting to %s, attempt %d of %d", c.host, i, c.retries)
		conn, err := grpcConnectWithTimeOut(ctx, c.host, c.timeout, c.grpcPingTime, c.timeout)
		if err != nil {
			c.log.Warnf("can't establish connection to %s, sleep for %s", c.host, c.timeout)
			SleepInContext(ctx, c.retryInterval)

			continue
		}

		return conn, nil
	}

	return nil, errors.Errorf("max connection attempts reached: %d", c.retries)
}

func grpcConnectWithTimeOut(ctx context.Context, target string, connectTimeOut, grpcPingTime, grpcPingTimeOut time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeOut)
	defer cancel()

	return grpc.DialContext(
		ctx,
		target,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    grpcPingTime,
			Timeout: grpcPingTimeOut,
		}))
}
