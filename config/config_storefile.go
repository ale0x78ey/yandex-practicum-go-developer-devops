package config

import (
	"flag"
)

const (
	DefaultInitStore = true
	DefaultStoreFile = "/tmp/devops-metrics-db.json"
)

type StoreFileConfig struct {
	InitStore bool   `env:"RESTORE"`
	StoreFile string `env:"STORE_FILE"`
}

func NewStoreFileConfig() *StoreFileConfig {
	cfg := StoreFileConfig{}
	flag.BoolVar(&cfg.InitStore, "r", DefaultInitStore, "RESTORE")
	flag.StringVar(&cfg.StoreFile, "f", DefaultStoreFile, "STORE_FILE")
	return &cfg
}
