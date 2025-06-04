package handler

import (
	"context"
	"log/slog"
	"url-shortener/internal/models"
)

// Service defines the interface for URL shortening operations
type Service interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(req []string, userID string) error
}

// Handler manages HTTP request handling for URL shortening service
type Handler struct {
	service Service
	log     *slog.Logger
}

// New creates a new Handler instance
func New(s Service, log *slog.Logger) *Handler {
	return &Handler{
		service: s,
		log:     log,
	}
}
