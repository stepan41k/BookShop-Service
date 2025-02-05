package postgres

import (
	"context"

	"github.com/stepan41k/testMidlware/internal/storage"
)

func (p *PGPool) CreateAuthor(newAuthor storage.Author) (id int, err error) {
	
	err = p.pool.QueryRow(context.Background(), `
		INSERT INTO authors (name)
		VALUES ($1)
		RETURNING id
	`, newAuthor.Name).Scan(
		id,
	)
	
	return id, err
}

func (p *PGPool) ReadAuthors() ([]storage.Author, error) {
	rows, err := p.pool.Query(context.Background(),`
	SELECT id, name
	FROM authors;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []storage.Author

	for rows.Next() {
		var item storage.Author
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

func (p *PGPool) ReadAuthorById(id int) (storage.Author, error) {
	var data storage.Author

	err := p.pool.QueryRow(context.Background(), `
		SELECT id, name
		FROM authors
		WHERE id = $1;
	`, id).Scan(
		&data.ID,
		&data.Name,
	)

	return data, err

}

func (p *PGPool) UpdateAuthor(author string, updAuthor storage.Author) (error) {
	
	_, err := p.pool.Exec(context.Background(), `
		UPDATE authors
		SET name = $1
		WHERE name = $2
	`, updAuthor.Name, author)

	return err
}

func (p *PGPool) DeleteAuthor(author string) (error) {

	_, err := p.pool.Exec(context.Background(), `
		DELETE FROM authors
		WHERE name = $1
	`, author)

	return err

}

