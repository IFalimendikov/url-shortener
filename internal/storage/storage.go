package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
)

// Storage holds the file, database and URL mapping information
type Storage struct {
	cfg  *config.Config
	DB   *sql.DB
	File os.File
	URLs map[string]models.URLRecord
}

// Query for creating urls table
var (
	UrlsQuery = `CREATE TABLE IF NOT EXISTS urls (user_id text, short_url text, url text PRIMARY KEY, deleted bool DEFAULT false);`
)

// New creates a new Storage instance with file and database connections
func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
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

		_, err = db.Exec(UrlsQuery)
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
