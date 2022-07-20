package utils

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ClientConnector struct {
	log           *zap.SugaredLogger
	host          string
	retries       int
	timeout       time.Duration
	retryInterval time.Duration
}

func NewClientConnector(log *zap.SugaredLogger,
	host string,
	retries int,
	timeout,
	retryInterval time.Duration) ClientConnector {
	return ClientConnector{
		log:           log,
		host:          host,
		retries:       retries,
		timeout:       timeout,
		retryInterval: retryInterval,
	}
}

func (c *ClientConnector) ConnectWithRetry(ctx context.Context) (net.Conn, error) {
	for i := 1; i <= c.retries; i++ {
		c.log.Debugf("connecting to %s, attempt %d of %d", c.host, i, c.retries)
		conn, err := connectWithTimeout(c.host, c.timeout)
		if err != nil {
			c.log.Warnf("can't establish connection to %s, sleep for %s", c.host, c.timeout)
			SleepInContext(ctx, c.retryInterval)

			continue
		}
		return conn, nil
	}
	return nil, errors.Errorf("max connection attempts reached %d", c.retries)
}

func connectWithTimeout(target string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
