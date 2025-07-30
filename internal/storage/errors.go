package storage

import (
	"errors"
)

// Package level errors for the URL shortener service layer
var (
	ErrorDuplicate  = errors.New("duplicate URL record")
	ErrorNotFound   = errors.New("error finding URL")
	ErrorURLDeleted = errors.New("URL was deleted")
	ErrorURLSave    = errors.New("can't save URL")
	ErrorTxCommit   = errors.New("can't commit a Tx")
)
