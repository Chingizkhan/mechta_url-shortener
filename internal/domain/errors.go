package domain

import "errors"

var (
	ErrAlreadyExists = errors.New("already_exists")
	ErrNotFound      = errors.New("not_found")
	ErrEmptyLink     = errors.New("empty_link")
	ErrTooShort      = errors.New("too_short")
	ErrInvalidURL    = errors.New("invalid_url")
)
