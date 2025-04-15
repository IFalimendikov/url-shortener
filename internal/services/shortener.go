package services

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"

	"github.com/deatil/go-encoding/base62"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorDatabase  *pgconn.PgError
	ErrorDuplicate = errors.New("duplicate URL record")
	ErrorNotFound  = errors.New("error finding URL")
)

type URLService interface {
	ServSave(ctx context.Context, url, userID string) (string, error)
	ServGet(shortURL string) (string, error)
	ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error
	GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error
	PingDB() bool
}

type URLStorage struct {
	Context context.Context
	MU      sync.RWMutex
	Log     *slog.Logger
	Storage *storage.Storage
	Encoder *json.Encoder
}

func NewURLService(ctx context.Context, log *slog.Logger, storage *storage.Storage) *URLStorage {
	service := &URLStorage{
		Context: ctx,
		Storage: storage,
		Log:     log,
		Encoder: json.NewEncoder(&storage.File),
	}
	return service
}

func (s *URLStorage) ServSave(ctx context.Context, url, userID string) (string, error) {
	short := base62.StdEncoding.EncodeToString([]byte(url))

	rec := storage.URLRecord{
		ID:       s.Storage.Count,
		UserID:   userID,
		ShortURL: short,
		URL:      url,
	}

	if s.Storage.DB != nil {
		_, err := s.Storage.DB.ExecContext(ctx, storage.SaveURL, rec.ID, rec.UserID, rec.ShortURL, rec.URL)
		if err != nil {
			if errors.As(err, &ErrorDatabase) && ErrorDatabase.Code == pgerrcode.UniqueViolation {
				return short, ErrorDuplicate
			}
			return "", errors.New("failed to save URL to the database")
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
		s.Storage.Count++
		s.MU.Unlock()
	}

	return short, nil
}

func (s *URLStorage) ServGet(shortURL string) (string, error) {
	if s.Storage.DB != nil {
		var url string
		row := s.Storage.DB.QueryRow(storage.GetURL, shortURL)

		err := row.Scan(&url)
		if err != nil {
			return "", nil
		}

		if url != "" {
			return url, nil
		}
		return "", ErrorNotFound
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

func (s *URLStorage) PingDB() bool {
	if s.Storage.DB != nil {
		err := s.Storage.DB.Ping()
		return err == nil
	}
	return false
}

func (s *URLStorage) ShortenBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest, res *[]models.BatchUnitURLResponse) error {
	db := s.Storage.DB
	if db != nil {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		stmt, err := db.PrepareContext(ctx, storage.SaveURL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, x := range req {
			id := s.Storage.Count
			short := base62.StdEncoding.EncodeToString([]byte(x.URL))

			_, err = stmt.ExecContext(ctx, id, userID, short, x.URL)
			if err != nil {
				return err
			}
		}
		tx.Commit()
	}

	for _, x := range req {
		rec := storage.URLRecord{
			ID:       s.Storage.Count,
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
			s.Storage.Count++
			s.MU.Unlock()
		}
	}

	return nil
}

func (s *URLStorage) GetUserURLs(ctx context.Context, userID string, res *[]models.UserURLResponse) error {
	db := s.Storage.DB
	if db != nil {
		stmt, err := db.PrepareContext(ctx, storage.GetUserURL)
		if err != nil {
			return err
		}
		defer stmt.Close()

		rows, err := stmt.QueryContext(ctx, userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var url models.UserURLResponse
			err := rows.Scan(&url.ShortURL, &url.OriginalURL)
			if err != nil {
				return err
			}
			*res = append(*res, url)
		}

		if len(*res) == 0 {
			return ErrorNotFound
		}
	}
	return nil
}