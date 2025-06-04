package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

// Config holds the application configuration settings.
// It can be populated from environment variables using the env tags.
type Config struct {
	// ServerAddr specifies the server address in format host:port
	ServerAddr string `env:"SERVER_ADDRESS"`
	// BaseURL is the base URL for the shortened URLs
	BaseURL string `env:"BASE_URL"`
	// StoragePath specifies the path to the file storage
	StoragePath string `env:"FILE_STORAGE_PATH"`
	// DBAddress holds the database connection string
	DBAddress string `env:"DATABASE_DSN" envDefault:""`
}

// Read parses environment variables into the Config struct.
// It sets default values for ServerAddr, BaseURL, and StoragePath if they are not provided.
// The function will log.Fatal if environment parsing fails.
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
