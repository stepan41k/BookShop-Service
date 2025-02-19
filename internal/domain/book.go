package domain

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name" validate:"required"`
	Author string  `json:"author" validate:"required"`
	Genre  string  `json:"genre" validate:"required"`
	Price  string `json:"price" validate:"required"`
}