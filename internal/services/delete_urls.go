package services

import (
	"context"
	"database/sql"
	"sync"
	"time"
	"url-shortener/internal/models"
)

func (s *URLs) DeleteURLs(req []string, userID string) error {
	ctx := context.Background()
	ch := make(chan models.DeleteRecord, len(req))
	defer close(ch)

	for _, x := range req {
		del := models.DeleteRecord{
			UserID:   userID,
			ShortURL: x,
		}
		ch <- del
	}

	return s.processURLs(ctx, ch)
}

func (s *URLs) processURLs(ctx context.Context, chs ...chan models.DeleteRecord) error {
	var wg sync.WaitGroup
	var buffer []models.DeleteRecord
	resultCh := make(chan models.DeleteRecord, 20)
	timer := time.NewTicker(1 * time.Second)

	s.Log.Info("Starting URL processing", "channels", len(chs))

	go func() {
		wg.Wait()
		s.Log.Debug("All goroutines completed, closing result channel")
		close(resultCh)
	}()

	for _, ch := range chs {
		wg.Add(1)
		cha := ch
		go func() {
			defer wg.Done()
			for data := range cha {
				s.Log.Debug("Processing URL", "shortURL", data.ShortURL, "userID", data.UserID)
				resultCh <- data
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			s.Log.Warn("Context cancelled in background",
				"error", ctx.Err(),
				"buffered_urls", len(buffer))
			return ctx.Err()

		case x, ok := <-resultCh:
			if !ok {
				if len(buffer) > 0 {
					s.Log.Info("Processing final batch", "count", len(buffer))
					return s.commitDB(ctx, buffer)
				}
				s.Log.Info("URL processing completed")
				return nil
			}
			buffer = append(buffer, x)
			if len(buffer) >= 10 {
				s.Log.Info("Buffer full, committing batch", "count", len(buffer))
				if err := s.commitDB(ctx, buffer); err != nil {
					s.Log.Error("Failed to commit batch",
						"error", err,
						"batch_size", len(buffer))
					return err
				}
				buffer = buffer[:0]
			}

		case <-timer.C:
			if len(buffer) > 0 {
				s.Log.Info("Timer triggered, committing remaining URLs",
					"count", len(buffer))
				return s.commitDB(ctx, buffer)
			}
		}
	}
}

func (s *URLs) commitDB(ctx context.Context, records []models.DeleteRecord) error {
	db := s.Storage.DB
	if db != nil {
		tx, err := s.Storage.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		if err = s.Storage.Delete(ctx, tx, records); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
