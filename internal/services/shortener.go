package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"url-shortener/internal/storage"
	"url-shortener/internal/models"

	"github.com/deatil/go-encoding/base62"
	"go.uber.org/zap"
)

type URLService interface {
	ServSave(url string) (string, error)
	ServGet(shortURL string) (string, error)
	ShortenBatch(ctx context.Context, req models.ShortenURLBatchRequest, res *models.ShortenURLBatchResponse) error
	PingDB() bool
}

type URLStorage struct {
	Context context.Context
	MU      sync.RWMutex
	Log     *zap.SugaredLogger
	Storage *storage.Storage
	Encoder *json.Encoder
	Base    *base62.Encoding
}

func NewURLService(ctx context.Context, log *zap.SugaredLogger, storage *storage.Storage) *URLStorage {
	service := &URLStorage{
		Context: ctx,
		Storage: storage,
		Log:     log,
		Encoder: json.NewEncoder(&storage.File),
		Base:    base62.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"),
	}
	return service
}

func (s *URLStorage) ServSave(url string) (string, error) {
	short := base62.StdEncoding.EncodeToString([]byte(url))

	rec := storage.URLRecord{
		ID:       s.Storage.Count,
		ShortURL: short,
		URL:      url,
	}

	if s.Storage.DB != nil {
		_, err := s.Storage.DB.Exec(storage.SaveURL, rec.ID, rec.ShortURL, rec.URL)
		if err != nil {
			return "", fmt.Errorf("failed to save URL to database: %w", err)
		}
	}

	_, ok := s.Storage.URLs[rec.ShortURL]
	if !ok {
		s.MU.Lock()
		err := s.Encoder.Encode(rec)
		if err != nil {
			return "", err
		}

		s.Storage.URLs[rec.ShortURL] = rec
		s.Storage.Count++
		s.MU.Unlock()
	}

	return short, nil
}

func (s *URLStorage) ServGet(shortURL string) (string, error) {
	if s.Storage.DB != nil {
		var url string
		row := s.Storage.DB.QueryRow(storage.GetURL, shortURL)

		err := row.Scan(&url)
		if err != nil {
			return "", nil
		}

		if url != "" {
			return url, nil
		}
		return "", fmt.Errorf("URL not found")
	}

	s.MU.RLock()
	url, ok := s.Storage.URLs[shortURL]
	s.MU.RUnlock()
	if ok {
		return url.URL, nil
	} else {
		return "", fmt.Errorf("URL not found")
	}
}

func (s *URLStorage) PingDB() bool {
	if s.Storage.DB != nil {
		err := s.Storage.DB.Ping()
		return err == nil
	}
	return false
}

func (s *URLStorage) ShortenBatch(ctx context.Context, req models.ShortenURLBatchRequest, res *models.ShortenURLBatchResponse) error {
	db := s.Storage.DB
	if db != nil {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		stmt, err := db.PrepareContext(ctx, storage.SaveURL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, x := range req.URLs {
			id := s.Storage.Count
			short := base62.StdEncoding.EncodeToString([]byte(x.URL))

			_, err = stmt.ExecContext(ctx, id, short, x.URL)
			if err != nil {
				return err
			}
		}
		tx.Commit()
	}

	for _, x := range req.URLs {
		rec := storage.URLRecord{
			ID:       s.Storage.Count,
			ShortURL: base62.StdEncoding.EncodeToString([]byte(x.URL)),
			URL:      x.URL,
		}

		res.URLs = append(res.URLs, models.BatchUnitURLResponse{
			ID:    x.ID,
			Short: rec.ShortURL,
		})

		_, ok := s.Storage.URLs[rec.ShortURL]
		if !ok {
			s.MU.Lock()
			err := s.Encoder.Encode(rec)
			if err != nil {
				return err
			}
			s.Storage.URLs[rec.ShortURL] = rec
			s.Storage.Count++
			s.MU.Unlock()
		}
	}

	return nil
}
