package storage

import (
	"context"
	"url-shortener/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Delete(ctx context.Context, runner sq.BaseRunner, records []models.DeleteRecord) error {
    for _, x := range records {
        _, err := sq.Update("urls").
            Set("deleted", true).
            Where(sq.And{
                sq.Eq{"user_id": x.UserID},
                sq.Eq{"short_url": x.ShortURL},
            }).
            PlaceholderFormat(sq.Dollar).
            RunWith(runner).
            ExecContext(ctx)
            
        if err != nil {
            return err
        }
    }
    return nil
}