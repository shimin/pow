package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host          string        `mapstructure:"CLIENT_AUTH_HOST"` //= "127.0.0.1:5000"
	TargetBits    uint16        `mapstructure:"TARGET_BITS"`      //= 24
	Retries       int           `mapstructure:"RETRIES"`          //= 2
	RetryInterval time.Duration `mapstructure:"RETRY_INTERVAL"`   //= 2 * time.Second
	ConnTimeout   time.Duration `mapstructure:"CONN_TIMEOUT"`     //= 2 * time.Second
	GrpcPingTime  time.Duration `mapstructure:"GRPC_PING_TIME"`   //= 30 * time.Second
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
