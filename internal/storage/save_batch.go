package storage

import (
	"context"
	"github.com/deatil/go-encoding/base62"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) SaveBatch(ctx context.Context, runner sq.BaseRunner, userID string, req []models.BatchUnitURLRequest) error {
	for _, x := range req {
		short := base62.StdEncoding.EncodeToString([]byte(x.URL))

		_, err := sq.Insert("urls").
			Columns("user_id", "short_url", "url").
			Values(userID, short, x.URL).
			RunWith(runner).
			PlaceholderFormat(sq.Dollar).
			ExecContext(ctx)

		if err != nil {
			return err
		}
	}
	return nil
}
