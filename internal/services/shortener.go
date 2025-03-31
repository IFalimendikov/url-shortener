package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"url-shortener/internal/storage"

	"github.com/deatil/go-encoding/base62"
	"go.uber.org/zap"
)

type URLService interface {
	ServSave(url string) (string, error)
	ServGet(shortURL []byte) (string, error)
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
	short := s.Base.EncodeToString([]byte(url))

	rec := storage.URLRecord{
		ID:       s.Storage.Count,
		ShortURL: short,
		URL:      url,
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
	err := s.Storage.DB.Ping() 

	if err != nil {
		var url string
		row := s.Storage.DB.QueryRow(storage.GetURL, shortURL)
		
		err := row.Scan(url)
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

func (s *URLStorage) PingDB () bool {
	err := s.Storage.DB.Ping()
	return err != nil
}