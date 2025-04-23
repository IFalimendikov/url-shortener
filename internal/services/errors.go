package services

import (
	"errors"
)

var (
	ErrorDuplicate  = errors.New("duplicate URL record")
	ErrorNotFound   = errors.New("error finding URL")
	ErrorURLDeleted = errors.New("URL was deleted")
	ErrorURlSave = errors.New("can't save URL")
)