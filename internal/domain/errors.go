package domain

import "errors"

var (
	ErrInvalidFilterRange = errors.New("invalid filter range")
	ErrBookNotFound       = errors.New("book not found")
)
