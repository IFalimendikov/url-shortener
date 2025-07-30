package logger

import (
	"log/slog"
	"os"
)

// NewLogger creates a new structured logger with the specified log level
func New() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return log
}
