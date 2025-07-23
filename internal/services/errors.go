package services

import (
	"errors"
)

// Package level errors for the URL shortener service layer
var (
	ErrorNotFound = errors.New("error finding URL")
	ErrorNoDB     = errors.New("error connecting DB")
)
