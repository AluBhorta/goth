package customerrors

import "errors"

var ErrNotFound = errors.New("not found")

var ErrDuplicateKey = errors.New("duplicate key")
