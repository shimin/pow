package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host         string        `envconfig:"SERVER_LISTEN_ADDRESS"` // = ":5000"
	KeySize      uint16        `envconfig:"KEY_SIZE"`              // = 40
	TargetBits   uint16        `envconfig:"TARGET_BITS"`           // = 24
	GrpcPingTime time.Duration `envconfig:"GRPC_PING_TIME"`        // = 30 * time.Second
	GrpcTimeOut  time.Duration `envconfig:"GRPC_TIMEOUT"`          // = 60 * time.Second
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	err = envconfig.Process("", &config)
	return
}
