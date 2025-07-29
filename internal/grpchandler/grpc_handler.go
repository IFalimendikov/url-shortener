package grpchandler

import (
	"context"
	"log/slog"
	"url-shortener/internal/models"
	"url-shortener/internal/proto"

	"url-shortener/internal/config"
)

// Service defines the interface for URL shortening operations
type Service interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(req []string, userID string) error
	GetStats(ctx context.Context) (models.Stats, error)
}

// Handler manages GRPC request handling for URL shortening service
type GRPCHandler struct {
	proto.UnimplementedURLShortenerServer
	service Service
	cfg     config.Config
	log     *slog.Logger
}

// New creates a new GRPCHandler instance
func New(ctx context.Context, s Service, cfg config.Config, log *slog.Logger) *GRPCHandler {
	return &GRPCHandler{
		service: s,
		cfg:     cfg,
		log:     log,
	}
}
