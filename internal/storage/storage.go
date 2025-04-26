package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	"database/sql"
	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/deatil/go-encoding/base62"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg  *config.Config
	DB   *sql.DB
	File os.File
	URLs map[string]models.URLRecord
}

func NewStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
	file, err := os.OpenFile(cfg.StoragePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]models.URLRecord)
	records := []models.URLRecord{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record models.URLRecord
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for _, record := range records {
		urls[record.URL] = record
	}

	var db *sql.DB
	if cfg.DBAddress != "" {
		db, err = sql.Open("pgx", cfg.DBAddress)
		if err != nil {
			return nil, err
		}
		var query = `CREATE TABLE IF NOT EXISTS urls (user_id text, short_url text, url text PRIMARY KEY, deleted bool DEFAULT false);`

		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	storage := Storage{
		cfg:  cfg,
		File: *file,
		URLs: urls,
		DB:   db,
	}

	return &storage, err
}

func (s *Storage) PingDB() bool {
	if s.DB != nil {
		err := s.DB.Ping()
		return err == nil
	}
	return false
}

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

func (s *Storage) Get(ctx context.Context, shortURL string) (string, error) {
	var url models.URLRecord
	var query = `SELECT url, deleted FROM urls WHERE short_url = $1`
	row := s.DB.QueryRowContext(ctx, query, shortURL)

	err := row.Scan(&url.URL, &url.Deleted)
	if err != nil {
		return "", err
	}

	if url.Deleted {
		return "", ErrorURLDeleted
	}

	if url.URL != "" {
		return url.URL, nil
	}
	return "", ErrorNotFound
}

func (s *Storage) SaveBatch(ctx context.Context, userID string, req []models.BatchUnitURLRequest) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var query = `INSERT into urls (user_id, short_url, url) VALUES ($1, $2, $3)`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, x := range req {
		short := base62.StdEncoding.EncodeToString([]byte(x.URL))

		_, err = stmt.ExecContext(ctx, userID, short, x.URL)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return ErrorTxCommit
	}
	return nil
}

func (s *Storage) GetMultiple(ctx context.Context, userID string, res *[]models.UserURLResponse) error {
	var query = `SELECT short_url, url FROM urls WHERE user_id = $1`
	stmt, err := s.DB.PrepareContext(ctx, query)
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
	return nil
}

func (s *Storage) Delete(ctx context.Context, records []models.DeleteRecord) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var query = `UPDATE urls SET deleted = true WHERE user_id = $1 AND short_url = $2`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, x := range records {
		_, err := stmt.ExecContext(ctx, x.UserID, x.ShortURL)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return ErrorTxCommit
	}
	return nil
}
