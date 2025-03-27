package flag

import (
	"flag"

	"url-shortener/internal/config"
)

func ParseFlags() config.Config {
	cfg := config.Config{}

	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server host address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base HTTP address returned before short URL")
	flag.StringVar(&cfg.StoragePath, "f", cfg.StoragePath, "Storage file path for URLs")
	flag.StringVar(&cfg.DBAddress, "d", cfg.DBAddress, "Database connection.")
	flag.Parse()

	return cfg
}
