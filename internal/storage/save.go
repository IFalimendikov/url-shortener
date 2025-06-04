package storage

import (
	"context"
	"errors"
	"url-shortener/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	// sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Save(ctx context.Context, rec models.URLRecord) error {
	var query = `INSERT into urls (user_id, short_url, url) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, query, rec.UserID, rec.ShortURL, rec.URL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorDuplicate
		}
		return ErrorURLSave
	}
	return nil
}