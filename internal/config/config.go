package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddr  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	StoragePath string `env:"FILE_STORAGE_PATH"`
	DBAddress   string `env:DATABASE_DSN`
}

func Read(cfg *Config) {

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerAddr == "" {
		cfg.ServerAddr = "localhost:8080"
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://" + cfg.ServerAddr
	}

	if cfg.StoragePath == "" {
		cfg.StoragePath = "urls.json"
	}
}
