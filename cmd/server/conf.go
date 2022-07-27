package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host         string        `mapstructure:"SERVER_LISTEN_ADDRESS"` // = ":5000"
	KeySize      uint16        `mapstructure:"KEY_SIZE"`              // = 40
	TargetBits   uint16        `mapstructure:"TARGET_BITS"`           // = 24
	GrpcPingTime time.Duration `mapstructure:"GRPC_PING_TIME"`        // = 30 * time.Second
	GrpcTimeOut  time.Duration `mapstructure:"GRPC_TIMEOUT"`          // = 60 * time.Second
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
