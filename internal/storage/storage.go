package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"url-shortener/internal/config"
)

type Storage struct {
	cfg   *config.Config
	File  os.File
	Count uint
	URLs  map[string]URLRecord
}

type URLRecord struct {
	ID       uint   `json:"uuid"`
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}

func NewStorage(cfg *config.Config) (*Storage, error) {
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

	storage := Storage{
		cfg:   cfg,
		File:  *file,
		Count: count,
		URLs:  urls,
	}

	return &storage, err
}
