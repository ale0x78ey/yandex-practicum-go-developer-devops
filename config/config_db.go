package config

import (
	"flag"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/db"
)

func NewDBConfig() *db.Config {
	cfg := &db.Config{}
	flag.StringVar(&cfg.DSN, "d", "", "DATABASE_DSN")
	return cfg
}
