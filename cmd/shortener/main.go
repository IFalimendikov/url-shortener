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
	grpcHandler "url-shortener/internal/grpchandler"
	grpcTransport "url-shortener/internal/grpctransport"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"

	"net"

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

	httpHandler := handler.New(s, log)
	grpcHandler := grpcHandler.New(ctx, s, cfg, log)

	httpErrCh := make(chan error, 1)
	grpcErrCh := make(chan error, 1)

	httpTransport := transport.New(cfg, httpHandler, log)
	r := transport.NewRouter(httpTransport)

	grpcServer := grpcTransport.New(grpcHandler, log)
	g := grpcTransport.NewRouter(grpcServer)

	go func() {
		if cfg.HTTPS {
			if err := r.RunTLS(cfg.ServerAddr, "cert.pem", "key.pem"); err != nil {
				httpErrCh <- fmt.Errorf("HTTP server failed: %w", err)
			}
		} else {
			if err := r.Run(cfg.ServerAddr); err != nil {
				httpErrCh <- fmt.Errorf("HTTP server failed: %w", err)
			}
		}
	}()

	go func() {
		grpcAddr := ":50051"
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			grpcErrCh <- err
			return
		}
		log.Info("Starting gRPC server", "address", grpcAddr)
		if err := g.Serve(lis); err != nil {
			grpcErrCh <- err
		}
	}()

	select {
	case err := <-httpErrCh:
		log.Error("HTTP server error", "error", err)
	case err := <-grpcErrCh:
		log.Error("gRPC server error", "error", err)
	case <-ctx.Done():
		log.Info("Servers shut down successfully")
	}
}
