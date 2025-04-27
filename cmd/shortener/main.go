package main

import (
	"context"
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"
	"url-shortener/internal/handler"
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
	h := handler.NewHandler(s, log)

	t := transport.NewTransport(cfg, h, log)
	r := transport.NewRouter(t)
	r.Run(cfg.ServerAddr)
}
