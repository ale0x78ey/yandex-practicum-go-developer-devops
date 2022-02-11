package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/agent"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

type Config struct {
	Http      *HttpConfig
	Server    *server.Config
	StoreFile *StoreFileConfig
	Agent     *agent.Config
}

func LoadAgentConfig() *Config {
	conf := &Config{
		Http:  NewHttpConfig(),
		Agent: NewAgentConfig(),
	}

	flag.Parse()

	if err := env.Parse(conf.Http); err != nil {
		log.Fatalf("Failed to parse http config options: %v", err)
	}

	if err := env.Parse(conf.Agent); err != nil {
		log.Fatalf("Failed to parse agent config options: %v", err)
	}

	return conf
}

func LoadServerConfig() *Config {
	conf := &Config{
		Http:      NewHttpConfig(),
		Server:    NewServerConfig(),
		StoreFile: NewStoreFileConfig(),
	}

	flag.Parse()

	if err := env.Parse(conf.Http); err != nil {
		log.Fatalf("Failed to parse http config options: %v", err)
	}

	if err := env.Parse(conf.Server); err != nil {
		log.Fatalf("Failed to parse server config options: %v", err)
	}

	if err := env.Parse(conf.StoreFile); err != nil {
		log.Fatalf("Failed to parse store file config options: %v", err)
	}

	return conf
}
