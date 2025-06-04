package flag

import (
    "flag"
    "url-shortener/internal/config"
)

// Parse reads and parses command line flags into a Config struct.
// It handles the following flags:
//   -a: HTTP server host address
//   -b: Base HTTP address returned before short URL
//   -f: Storage file path for URLs
//   -d: Database connection string
// Returns a populated Config struct with the parsed values.
func Parse() config.Config {
    cfg := config.Config{}

    flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server host address")
    flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base HTTP address returned before short URL")
    flag.StringVar(&cfg.StoragePath, "f", cfg.StoragePath, "Storage file path for URLs")
    flag.StringVar(&cfg.DBAddress, "d", cfg.DBAddress, "Database connection.")
    flag.Parse()

    return cfg
}