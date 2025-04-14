package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"url-shortener/internal/config"
)

type Storage struct {
	cfg   *config.Config
	DB    *sql.DB
	File  os.File
	Count uint
	URLs  map[string]URLRecord
}

type URLRecord struct {
	ID       uint   `json:"uuid"`
	UserId   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
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

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	urls := make(map[string]URLRecord)
	records := []URLRecord{}

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
	}

	storage := Storage{
		cfg:   cfg,
		File:  *file,
		Count: count,
		URLs:  urls,
		DB:    db,
	}

	return &storage, err
}
