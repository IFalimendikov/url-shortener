package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/services"
	"url-shortener/internal/transport"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage"
)

func main() {
	cfg := flag.ParseFlags()
	config.Read(&cfg)
	log := logger.NewLogger()

	store, err := storage.NewStorage(&cfg)
	if err != nil {
		log.Fatalf("Error creating new storage:%s", err)
	}
	s := services.NewURLService(log, store)

	t := transport.NewTransport(cfg, s, log)
	r := transport.NewRouter(cfg, t)
	r.Run(cfg.ServerAddr)
}
