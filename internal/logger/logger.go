package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return log
}