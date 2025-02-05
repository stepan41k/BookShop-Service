package storage

import (
	"errors"
)

type Book struct {
	ID			int 	`json:"id"`
	Name		string 	`json:"name"`
	AuthorID	int 	`json:"author_id"`
	GenreID 	int 	`json:"genre_id"`
	Price 		int 	`json:"price"`
}

type Genre struct {
	ID int				`json:"id"`
	Genre string		`json:"genre"`
}

type Author struct {
	ID int				`json:"id"`
	Name string			`json:"name"`
}

var (
	ErrBookNotFound = errors.New("book not found")
	ErrAuthorNotFound = errors.New("author not found")
	ErrGenreNotFound = errors.New("genre not found")

	ErrBookExists = errors.New("book already exists")
	ErrAuthorExists = errors.New("author already exists")
	ErrGenreExists = errors.New("genre already exists")
)