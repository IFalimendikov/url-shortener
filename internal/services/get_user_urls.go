package services

import (
	"context"
	"url-shortener/internal/models"
)

func (s *URLs) GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error {
	db := s.Storage.DB
	if db != nil {
		err := s.Storage.GetMultiple(ctx, userID, res)
		if err != nil {
			return err
		}
	}
	return nil
}
