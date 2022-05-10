package config

import (
	"flag"
)

const (
	DefaultServerAddress = "127.0.0.1:8080"
)

type HTTPConfig struct {
	ServerAddress string `env:"ADDRESS"`
}

func NewHTTPConfig() *HTTPConfig {
	cfg := HTTPConfig{}
	flag.StringVar(&cfg.ServerAddress, "a", DefaultServerAddress, "ADDRESS")
	return &cfg
}
