package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return log
}