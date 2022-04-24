package config

import (
	"flag"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/db"
)

const (
	DefaultMigrationsURL = "file://migrations"
)

func NewDBConfig() *db.Config {
	cfg := &db.Config{
		MigrationsURL: DefaultMigrationsURL,
	}
	flag.StringVar(&cfg.DSN, "d", "", "DATABASE_DSN")
	return cfg
}
