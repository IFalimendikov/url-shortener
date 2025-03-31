package services

import (
	"encoding/json"
	"errors"
	"sync"
	"url-shortener/internal/storage"

	"github.com/deatil/go-encoding/base62"
	"go.uber.org/zap"
)

var (
	ErrURLNotFound = errors.New("URL not found")
)

type URLService interface {
	ServSave(url string) (string, error)
	ServGet(shortURL string) (string, error)
}

type URLStorage struct {
	MU      sync.RWMutex
	Log     *zap.SugaredLogger
	Storage *storage.Storage
	Encoder *json.Encoder
}

func NewURLService(log *zap.SugaredLogger, storage *storage.Storage) *URLStorage {
	service := &URLStorage{
		Storage: storage,
		Log:     log,
		Encoder: json.NewEncoder(&storage.File),
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
	s.MU.RLock()
	url, ok := s.Storage.URLs[shortURL]
	s.MU.RUnlock()
	if ok {
		return url.URL, nil
	} else {
		return "", ErrURLNotFound
	}
}
