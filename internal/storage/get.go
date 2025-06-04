package storage

import (
	"context"
	"url-shortener/internal/models"

	// sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	var url models.URLRecord
	var query = `SELECT url, deleted FROM urls WHERE short_url = $1`
	row := s.DB.QueryRowContext(ctx, query, shortURL)

	err := row.Scan(&url.URL, &url.Deleted)
	if err != nil {
		return "", err
	}

	if url.Deleted {
		return "", ErrorURLDeleted
	}

	if url.URL != "" {
		return url.URL, nil
	}
	return "", ErrorNotFound
}
