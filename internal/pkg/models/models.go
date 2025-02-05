package models

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