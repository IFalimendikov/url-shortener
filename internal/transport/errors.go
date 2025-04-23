package transport

import (
	"errors"
)

var (
	ErrorDuplicate = errors.New("duplicate URL record")
	ErrorNotFound  = errors.New("error finding URL")
)