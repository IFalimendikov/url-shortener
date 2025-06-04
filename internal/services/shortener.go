package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
)

// Service defines the interface for URL shortening operations
type Service interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(ctx context.Context, req []string, userID string) error
}

// URLs implements the Service interface and manages URL shortening operations
type URLs struct {
    MU      sync.RWMutex      // Mutex for thread-safe operations
    Log     *slog.Logger      // Logger for service operations
    Storage *storage.Storage  // Storage interface for persistence
    Encoder *json.Encoder     // JSON encoder for data serialization
}

// New creates and initializes a new URLs service instance
func New(ctx context.Context, log *slog.Logger, storage *storage.Storage) *URLs {
    service := &URLs{
        Storage: storage,
        Log:     log,
        Encoder: json.NewEncoder(&storage.File),
    }
    return service
}