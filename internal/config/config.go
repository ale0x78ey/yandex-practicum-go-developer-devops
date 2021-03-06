package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/service/agent"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage/db"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage/file"
)

type Config struct {
	HTTP      *HTTPConfig
	Server    *server.Config
	StoreFile *file.Config
	Agent     *agent.Config
	DB        *db.Config
}

func LoadAgentConfig() *Config {
	conf := &Config{
		HTTP:  NewHTTPConfig(),
		Agent: NewAgentConfig(),
	}

	flag.Parse()

	if err := env.Parse(conf.HTTP); err != nil {
		log.Fatalf("Failed to parse http config options: %v", err)
	}

	if err := env.Parse(conf.Agent); err != nil {
		log.Fatalf("Failed to parse agent config options: %v", err)
	}

	return conf
}

func LoadServerConfig() *Config {
	conf := &Config{
		HTTP:      NewHTTPConfig(),
		Server:    NewServerConfig(),
		StoreFile: NewStoreFileConfig(),
		DB:        NewDBConfig(),
	}

	flag.Parse()

	if err := env.Parse(conf.HTTP); err != nil {
		log.Fatalf("Failed to parse http config options: %v", err)
	}

	if err := env.Parse(conf.Server); err != nil {
		log.Fatalf("Failed to parse server config options: %v", err)
	}

	if err := env.Parse(conf.StoreFile); err != nil {
		log.Fatalf("Failed to parse store file config options: %v", err)
	}

	if err := env.Parse(conf.DB); err != nil {
		log.Fatalf("Failed to parse DB config options: %v", err)
	}

	return conf
}
