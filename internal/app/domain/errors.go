package domain

import "errors"

var (
	ErrRequired = errors.New("required value")
	ErrNotFound = errors.New("not found")
)
