package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
)

type Service interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(ctx context.Context, req []string, userID string) error
}

type URLs struct {
	MU      sync.RWMutex
	Log     *slog.Logger
	Storage *storage.Storage
	Encoder *json.Encoder
}

func New(ctx context.Context, log *slog.Logger, storage *storage.Storage) *URLs {
	service := &URLs{
		Storage: storage,
		Log:     log,
		Encoder: json.NewEncoder(&storage.File),
	}
	return service
}
