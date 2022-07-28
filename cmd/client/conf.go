package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host          string        `envconfig:"CLIENT_AUTH_HOST"` //= "127.0.0.1:5000"
	TargetBits    uint16        `envconfig:"TARGET_BITS"`      //= 24
	Retries       int           `envconfig:"RETRIES"`          //= 2
	RetryInterval time.Duration `envconfig:"RETRY_INTERVAL"`   //= 2 * time.Second
	ConnTimeout   time.Duration `envconfig:"CONN_TIMEOUT"`     //= 2 * time.Second
	GrpcPingTime  time.Duration `envconfig:"GRPC_PING_TIME"`   //= 30 * time.Second
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	err = envconfig.Process("", &config)
	return
}
