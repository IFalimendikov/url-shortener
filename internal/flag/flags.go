package flag

import (
	"flag"

	"url-shortener/internal/config"
)

func ParseFlags() config.Config {
	config := config.Config{}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "HTTP server host address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base HTTP address returned before short URL")
	flag.StringVar(&config.StoragePath, "f", config.StoragePath, "Storage file path for URLs")
	flag.Parse()

	return config
}
