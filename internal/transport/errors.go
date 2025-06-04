package transport

import (
	"errors"
)

// Package level errors for the URL shortener transport layer
var (
	ErrorDuplicate = errors.New("duplicate URL record")
)
