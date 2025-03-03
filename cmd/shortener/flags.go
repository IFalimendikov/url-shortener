package main

import (
	"flag"
	
	"url-shortener/internal/app/config"
)

func ParseFlags() config.Config{
	config := config.Config{}
	
	flag.StringVar(&config.HTTPAddr, "a", config.HTTPAddr, "HTTP server host address")
	flag.StringVar(&config.BaseAddr, "b", config.BaseAddr, "Base HTTP address returned before short URL")
	flag.Parse()

	return config
}