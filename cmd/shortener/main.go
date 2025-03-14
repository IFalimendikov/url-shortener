package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/services"
	"url-shortener/internal/transport"
	"url-shortener/internal/logger"
)

func main() {
	cfg := flag.ParseFlags()
	config.Read(&cfg)
	log := logger.NewLogger()

	s := services.NewURLService(log)

	t := transport.NewTransport(cfg, s, log)
	r := transport.NewRouter(cfg, t)
	r.Run(cfg.ServerAddr)
}
