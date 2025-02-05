package repository

import (
	"context"
	"github.com/stepan41k/testMidlware/internal/pkg/models"
)

func (repo *PGRepo) GetBooks() ([]models.Book, error) {
	rows, err := repo.pool.Query(context.Background(), `
		SELECT id, name, author_id, genre_id, price
		FROM books;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []models.Book
	for rows.Next() {
		var item models.Book
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.AuthorID,
			&item.GenreID,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (repo *PGRepo) NewBook(item models.Book) (id int, err error) {
	err = repo.pool.QueryRow(context.Background(), `
		INSERT INTO books (name, author_id, genre_id, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`,
		item.Name,
		item.AuthorID,
		item.GenreID,
		item.Price,
		).Scan(&id)
	
	return id, err
}

func (repo *PGRepo) GetBookByID(id int) (models.Book, error) {
	var book models.Book
	err := repo.pool.QueryRow(context.Background(), `
		SELECT id, name, author_id, genre_id, price
		FROM books
		WHERE id = $1;
	`,	id).Scan(
		&book.ID,
		&book.Name,
		&book.AuthorID,
		&book.GenreID,
		&book.Price,
	)
	return book, err
}

func (repo *PGRepo) UpdateBook(item string, newBook models.Book)(models.Book, error) {

	var book models.Book

	_, err := repo.pool.Exec(context.Background(), `
	UPDATE books
	SET name = $1, author_id = $2, genre_id = $3, price = $4
	WHERE name = $5`,
	newBook.Name, newBook.AuthorID, newBook.GenreID, newBook.Price, item)
		
	return book, err
}

func (repo * PGRepo) DeleteBook(name string) (error) {
	_, err := repo.pool.Exec(context.Background(), `
		DELETE FROM books
		WHERE name = $1
		`, name)
	return err
}