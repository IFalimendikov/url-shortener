package storage

import (
	"context"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	var url models.URLRecord

	row := sq.Select("url", "deleted").
		From("urls").
		Where(sq.Eq{"short_url": shortURL}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.DB).
		QueryRowContext(ctx)

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
