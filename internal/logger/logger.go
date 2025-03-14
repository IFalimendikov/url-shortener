package logger

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.SugaredLogger) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()
	return sugar
}