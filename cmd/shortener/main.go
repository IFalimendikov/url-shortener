package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"

	"url-shortener/cmd/staticlint"
	_ "url-shortener/docs"
	"url-shortener/internal/config"
	"url-shortener/internal/flag"
	"url-shortener/internal/handler"
	"url-shortener/internal/logger"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"
	"url-shortener/internal/transport"

	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// @title           URL Shortener API
// @version         1.0
// @description     A URL shortening service API

// @host      localhost:8080
// @BasePath  /api/v1

const Config = `config.json`

type ConfigData struct {
	Staticcheck []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfgLint ConfigData
	if err = json.Unmarshal(data, &cfgLint); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		staticlint.ExitCheckAnalyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfgLint.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)

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
