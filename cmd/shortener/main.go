package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"

	_ "url-shortener/docs"
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// @title           URL Shortener API
// @version         1.0
// @description     A URL shortening service API

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	cfg := flag.Parse()
	config.New(&cfg)

	log := logger.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

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

	go func() {
		if cfg.HTTPS {
			r.RunTLS(cfg.ServerAddr, "cert.pem", "key.pem")
		} else {
			r.Run(cfg.ServerAddr)
		}
	}()

	<-ctx.Done()
	stop()
	log.Info("Received shutdown signal, shutting down gracefully...")
}
