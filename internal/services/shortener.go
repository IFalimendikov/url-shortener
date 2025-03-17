package services

import (
	"fmt"
	"sync"

	base "github.com/jcoene/go-base62"
	"go.uber.org/zap"
)

type URLService interface {
	ServShort(url string) (string, error)
	ServSave(url string) (string, error)
	ServGet(shortURL string) (string, error)
}

type URLStorage struct {
	urls    map[string]string
	counter int
	mu      sync.RWMutex
	log     zap.SugaredLogger
}

func NewURLService(log *zap.SugaredLogger) *URLStorage {
	service := &URLStorage{
		urls: make(map[string]string),
		log:  *log,
	}
	return service
}

func (s *URLStorage) ServSave(url string) (string, error) {
	s.mu.Lock()
	s.counter++
	urlShort := base.Encode(int64(s.counter))
	s.urls[urlShort] = url
	s.mu.Unlock()

	return urlShort, nil
}

func (s *URLStorage) ServGet(shortURL string) (string, error) {
	s.mu.RLock()
	url, ok := s.urls[shortURL]
	s.mu.RUnlock()
	if ok {
		return url, nil
	} else {
		return "", fmt.Errorf("URL not found")
	}
}
