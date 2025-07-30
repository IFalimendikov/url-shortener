package services

import (
	"context"
	"url-shortener/internal/models"

	"github.com/deatil/go-encoding/base62"
)

// SaveURL creates a shortened URL from the original URL and stores it with the associated userID
func (s *URLs) SaveURL(ctx context.Context, url, userID string) (string, error) {
	short := base62.StdEncoding.EncodeToString([]byte(url))

	rec := models.URLRecord{
		UserID:   userID,
		ShortURL: short,
		URL:      url,
	}

	if s.Storage.DB != nil {
		err := s.Storage.Save(ctx, rec)
		if err != nil {
			return rec.ShortURL, err
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
		s.MU.Unlock()
	}
	return short, nil
}
