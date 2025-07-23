package config

import (
	"encoding/json"
	"log"
	"os"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds the application configuration settings.
// It can be populated from environment variables using the env tags.
type Config struct {
	// ServerAddr specifies the server address in format host:port
	ServerAddr string `env:"SERVER_ADDRESS"`
	// BaseURL is the base URL for the shortened URLs
	BaseURL string `env:"BASE_URL" `
	// StoragePath specifies the path to the file storage
	StoragePath string `env:"FILE_STORAGE_PATH"`
	// DBAddress holds the database connection string
	DBAddress string `env:"DATABASE_DSN" envDefault:""`
	// Config in JSON format
	Config string `env:"CONFIG" envDefault:""`
	// HTTPS indicates whether the server should run with HTTPS
	HTTPS bool `env:"HTTPS"`
	// TrustedSubnet shows trusted subnets mask
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

type tempCfg struct {
	// ServerAddr specifies the server address in format host:port
	ServerAddr string `env:"server_address"`
	// BaseURL is the base URL for the shortened URLs
	BaseURL string `env:"base_url"`
	// StoragePath specifies the path to the file storage
	StoragePath string `env:"file_storage_path"`
	// DBAddress holds the database connection string
	DBAddress string `env:"database_dsn" envDefault:""`
	// HTTPS indicates whether the server should run with HTTPS
	HTTPS bool `env:"enable_https"`
	// TrustedSubnet shows trusted subnets mask
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
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

// New parses JSON variables into the Config struct.
func New(cfg *Config) error {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var tempCfg tempCfg
	if cfg.Config != "" {
		file, err := os.ReadFile(cfg.Config)
		if err != nil {
			return err
		}

		err = json.Unmarshal(file, &tempCfg)
		if err != nil {
			return err
		}

		if tempCfg.BaseURL != "" {
			cfg.BaseURL = tempCfg.BaseURL
		}

		if tempCfg.ServerAddr != "" {
			cfg.ServerAddr = tempCfg.ServerAddr
		}

		if tempCfg.StoragePath != "" {
			cfg.StoragePath = tempCfg.StoragePath
		}

		if tempCfg.TrustedSubnet != "" {
			cfg.TrustedSubnet = tempCfg.TrustedSubnet
		}

		if tempCfg.DBAddress != "" {
			cfg.DBAddress = tempCfg.DBAddress
		}
		// } else {
		// 	cfg.DBAddress = os.Getenv("DATABASE_DSN")
		// }

		if tempCfg.HTTPS {
			cfg.HTTPS = tempCfg.HTTPS
		}
	}
	Read(cfg)
	return nil
}
