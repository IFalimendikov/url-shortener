package main

import (
	"url-shortener/internal/app/config"
	"url-shortener/internal/app/transport"
)

func main() {
	cfg := ParseFlags()

	config.Read(&cfg)

	r := transport.NewURLRouter(cfg)
	r.Run(cfg.ServerAddr)
}
