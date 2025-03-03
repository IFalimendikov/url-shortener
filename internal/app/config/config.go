package config

import ()

type Config struct {
	HTTPAddr string
	BaseAddr string
}

func Read(cfg *Config) {

	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = "localhost:8080"
	}

	if cfg.BaseAddr == "" {
		cfg.BaseAddr = "http://" + cfg.HTTPAddr
	}
}
