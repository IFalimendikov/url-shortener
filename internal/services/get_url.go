package services

import (
	"context"
)

func (s *URLs) GetURL(ctx context.Context, shortURL string) (string, error) {
	if s.Storage.DB != nil {
		url, err := s.Storage.Get(ctx, shortURL)
		if err != nil {
			return url, err
		}
	}

	s.MU.RLock()
	url, ok := s.Storage.URLs[shortURL]
	s.MU.RUnlock()
	if ok {
		return url.URL, nil
	} else {
		return "", ErrorNotFound
	}
}
