package repository

import (
	"context"
	"github.com/stepan41k/testMidlware/internal/pkg/models"
)

func (repo *PGRepo) CreateAuthor(newAuthor models.Author) (id int, err error) {
	
	err = repo.pool.QueryRow(context.Background(), `
		INSERT INTO authors (name)
		VALUES ($1)
		RETURNING id
	`, newAuthor.Name).Scan(
		id,
	)
	
	return id, err
}

func (repo *PGRepo) ReadAuthors() ([]models.Author, error) {
	rows, err := repo.pool.Query(context.Background(),`
	SELECT id, name
	FROM authors;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []models.Author

	for rows.Next() {
		var item models.Author
		err = rows.Scan(
			&item.ID,
			&item.Name,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

func (repo *PGRepo) ReadAuthorById(id int) (models.Author, error) {
	var data models.Author

	err := repo.pool.QueryRow(context.Background(), `
		SELECT id, name
		FROM authors
		WHERE id = $1;
	`, id).Scan(
		&data.ID,
		&data.Name,
	)

	return data, err

}

func (repo *PGRepo) UpdateAuthor(author string, updAuthor models.Author) (error) {
	
	_, err := repo.pool.Exec(context.Background(), `
		UPDATE authors
		SET name = $1
		WHERE name = $2
	`, updAuthor.Name, author)

	return err
}

func (repo *PGRepo) DeleteAuthor(author string) (error) {

	_, err := repo.pool.Exec(context.Background(), `
		DELETE FROM authors
		WHERE name = $1
	`, author)

	return err

}

