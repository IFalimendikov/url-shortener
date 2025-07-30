package storage

import (
	"context"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
)

// Stats gets statistics of users from the database
func (s *Storage) Stats(ctx context.Context) (models.Stats, error) {
	var res models.Stats

	var urls int
	urlRow := sq.Select("COUNT(DISTINCT url)").
		From("urls").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.DB).
		QueryRowContext(ctx)

	err := urlRow.Scan(&urls)
	if err != nil {
		return res, err
	}

	var users int
	userRow := sq.Select("COUNT(DISTINCT user_id)").
		From("urls").
		PlaceholderFormat(sq.Dollar).
		RunWith(s.DB).
		QueryRowContext(ctx)

	err = userRow.Scan(&users)
	if err != nil {
		return res, err
	}

	res.Urls = urls
	res.Users = users

	return res, nil
}
