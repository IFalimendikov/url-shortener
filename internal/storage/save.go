package storage

import (
	"context"
	"errors"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Save stores a URL with its shortened version and user ID
func (s *Storage) Save(ctx context.Context, rec models.URLRecord) error {
	_, err := sq.Insert("urls").
		Columns("user_id", "short_url", "url").
		Values(rec.UserID, rec.ShortURL, rec.URL).
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorDuplicate
		}
		return ErrorURLSave
	}
	return nil
}
