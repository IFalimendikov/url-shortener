package services

import (
	"context"
	"database/sql"
	"url-shortener/internal/models"

	"github.com/deatil/go-encoding/base62"
)

func (s *URLs) ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error {
	if s.Storage.DB != nil {
		tx, err := s.Storage.DB.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if err = s.Storage.SaveBatch(ctx, tx, userID, req); err != nil {
			return err
		}
		return nil
	}

	for _, x := range req {
		rec := models.URLRecord{
			ShortURL: base62.StdEncoding.EncodeToString([]byte(x.URL)),
			URL:      x.URL,
		}

		*res = append(*res, models.BatchUnitURLResponse{
			ID:    x.ID,
			Short: rec.ShortURL,
		})

		_, ok := s.Storage.URLs[rec.ShortURL]
		if !ok {
			s.MU.Lock()
			err := s.Encoder.Encode(rec)
			if err != nil {
				return err
			}
			s.Storage.URLs[rec.ShortURL] = rec
			s.MU.Unlock()
		}
	}
	return nil
}
