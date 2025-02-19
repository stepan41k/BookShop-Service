package storage

import (
	"errors"
)

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrAuthorNotFound = errors.New("author not found")
	ErrGenreNotFound  = errors.New("genre not found")

	ErrBookExists   = errors.New("book already exists")
	ErrAuthorExists = errors.New("author already exists")
	ErrGenreExists  = errors.New("genre already exists")
)
