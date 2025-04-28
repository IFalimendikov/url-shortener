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
	cfg := flag.Parse()
	config.Read(&cfg)

	log := logger.New()
	ctx := context.Background()

	store, err := storage.New(ctx, &cfg)
	if err != nil {
		log.Error("Error creating new storage", "error", err)
	}
	defer store.File.Close()
	defer store.DB.Close()

	s := services.New(ctx, log, store)
	h := handler.New(s, log)

	t := transport.New(cfg, h, log)
	r := transport.NewRouter(t)
	r.Run(cfg.ServerAddr)
}
