package main

import (
	"flag"

	"url-shortener/internal/app/config"
)

func ParseFlags() config.Config {
	config := config.Config{}

	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "HTTP server host address")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "Base HTTP address returned before short URL")
	flag.Parse()

	return config
}
