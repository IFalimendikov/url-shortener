package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"database/sql"
	"url-shortener/internal/config"
	"url-shortener/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg   *config.Config
	DB    *sql.DB
	File  os.File

	URLs  map[string]types.URLRecord
}

func NewStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
	file, err := os.OpenFile(cfg.StoragePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	var count uint
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		count++
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]types.URLRecord)
	records := []types.URLRecord{}

	if count > 0 {
		dec := json.NewDecoder(file)
		dec.Decode(&records)

		for _, record := range records {
			urls[record.URL] = record
		}
	}

	var db *sql.DB
	if cfg.DBAddress != "" {
		db, err = sql.Open("pgx", cfg.DBAddress)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(CreateShortURLTable)
		if err != nil {
			return nil, err
		}
	}

	storage := Storage{
		cfg:   cfg,
		File:  *file,
		URLs:  urls,
		DB:    db,
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