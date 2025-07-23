package services

import (
	"context"
	"url-shortener/internal/models"
)

// GetStats returns user service statistics
func (s *URLs) GetStats(ctx context.Context) (models.Stats, error) {
	var res models.Stats

	if s.Storage.DB == nil {
		return res, ErrorNoDB
	}

	stats, err := s.Storage.Stats(ctx)
	if err != nil {
		return res, err
	}

	res.Urls = stats.Urls
	res.Users = stats.Users

	return res, nil
}
