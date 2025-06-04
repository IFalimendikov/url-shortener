package storage

import (
	"context"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetMultiple(ctx context.Context, userID string, res *[]models.UserURLResponse) error {
	rows, err := sq.Select("short_url", "url").
		From("urls").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.DB).
		QueryContext(ctx)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var url models.UserURLResponse
		err := rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			return err
		}
		*res = append(*res, url)
	}

	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
