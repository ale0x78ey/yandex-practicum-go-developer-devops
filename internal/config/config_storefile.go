package config

import (
	"flag"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage/file"
)

const (
	DefaultInitStore = true
	DefaultStoreFile = "/tmp/devops-metrics-db.json"
)

func NewStoreFileConfig() *file.Config {
	cfg := file.Config{}
	flag.BoolVar(&cfg.InitStore, "r", DefaultInitStore, "RESTORE")
	flag.StringVar(&cfg.StoreFile, "f", DefaultStoreFile, "STORE_FILE")
	return &cfg
}
