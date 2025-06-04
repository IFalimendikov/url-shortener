package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := flag.Parse()
	config.Read(&cfg)

	log := logger.New()
	ctx := context.Background()

	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Error("pprof server failed", "error", err)
		}
	}()

	if err := os.MkdirAll("profiles", 0755); err != nil {
		log.Error("Failed to create profiles directory", "error", err)
	}

	f, err := os.Create(filepath.Join("profiles", "base.pprof"))
	if err != nil {
		log.Error("Failed to create base profile", "error", err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Error("Failed to write heap profile", "error", err)
	}

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
