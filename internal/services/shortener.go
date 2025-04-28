package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/deatil/go-encoding/base62"
)

type URLRepository interface {
	SaveURL(ctx context.Context, url, userID string) (string, error)
	GetURL(ctx context.Context, shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
	DeleteURLs(ctx context.Context, req []string, userID string) error
}

type URLService struct {
	MU      sync.RWMutex
	Log     *slog.Logger
	Storage *storage.Storage
	Encoder *json.Encoder
}

func NewURLService(ctx context.Context, log *slog.Logger, storage *storage.Storage) *URLService {
	service := &URLService{
		Storage: storage,
		Log:     log,
		Encoder: json.NewEncoder(&storage.File),
	}
	return service
}

func (s *URLService) SaveURL(ctx context.Context, url, userID string) (string, error) {
	short := base62.StdEncoding.EncodeToString([]byte(url))

	rec := models.URLRecord{
		UserID:   userID,
		ShortURL: short,
		URL:      url,
	}

	if s.Storage.DB != nil {
		err := s.Storage.Save(ctx, rec)
		if err != nil {
			return rec.ShortURL, err
		}
	}

	_, ok := s.Storage.URLs[rec.ShortURL]
	if !ok {
		s.MU.Lock()
		err := s.Encoder.Encode(rec)
		if err != nil {
			return "", err
		}

		s.Storage.URLs[rec.ShortURL] = rec
		s.MU.Unlock()
	}
	return short, nil
}

func (s *URLService) GetURL(ctx context.Context, shortURL string) (string, error) {
	if s.Storage.DB != nil {
		url, err := s.Storage.Get(ctx, shortURL)
		if err != nil {
			return url, err
		}
	}

	s.MU.RLock()
	url, ok := s.Storage.URLs[shortURL]
	s.MU.RUnlock()
	if ok {
		return url.URL, nil
	} else {
		return "", ErrorNotFound
	}
}

func (s *URLService) PingDB() bool {
	return s.Storage.PingDB()
}

func (s *URLService) ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error {
	if s.Storage.DB != nil {
		err := s.Storage.SaveBatch(ctx, userID, req)
		if err != nil {
			return err
		}
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

func (s *URLService) GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error {
	db := s.Storage.DB
	if db != nil {
		err := s.Storage.GetMultiple(ctx, userID, res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *URLService) DeleteURLs(req []string, userID string) error {
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

func (s *URLService) processURLs(ctx context.Context, chs ...chan models.DeleteRecord) error {
    var wg sync.WaitGroup
    var buffer []models.DeleteRecord
    resultCh := make(chan models.DeleteRecord, 20)
    timer := time.NewTicker(5 * time.Second)

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

func (s *URLService) commitDB(ctx context.Context, records []models.DeleteRecord) error {
	db := s.Storage.DB
	if db != nil {
		err := s.Storage.Delete(ctx, records)
		if err != nil {
			return err
		}
	}
	return nil
}