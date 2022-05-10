package config

import (
	"flag"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/service/server"
)

func NewServerConfig() *server.Config {
	cfg := server.Config{
		ShutdownTimeout: server.DefaultShutdownTimeout,
	}
	flag.DurationVar(&cfg.StoreInterval, "i", server.DefaultStoreInterval, "STORE_INTERVAL")
	flag.StringVar(&cfg.Key, "k", "", "KEY")
	return &cfg
}
