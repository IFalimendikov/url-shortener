package main

import (
	"context"
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"
)

func main() {
	cfg := flag.ParseFlags()
	config.Read(&cfg)

	log := logger.NewLogger()
	ctx := context.Background()

	store, err := storage.NewStorage(ctx, &cfg)
	if err != nil {
		log.Error("Error creating new storage", "error", err)
	}
	defer store.File.Close()
	defer store.DB.Close()

	s := services.NewURLService(ctx, log, store)

	t := transport.NewTransport(cfg, s, log)
	r := transport.NewRouter(cfg, t)
	r.Run(cfg.ServerAddr)
}
