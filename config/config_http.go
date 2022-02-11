package config

import (
	"flag"
)

const (
	DefaultServerAddress = "127.0.0.1:8080"
)

type HttpConfig struct {
	ServerAddress string `env:"ADDRESS"`
}

func NewHttpConfig() *HttpConfig {
	cfg := HttpConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", DefaultServerAddress, "ADDRESS")
	return &cfg
}
