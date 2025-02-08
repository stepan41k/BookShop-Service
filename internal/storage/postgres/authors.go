package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *PGPool) SaveAuthor(newAuthor string) (int64, error) {
	
	const op = "storage.postgres.authors.SaveAuthor"

	row := p.pool.QueryRow(context.Background(), `
		INSERT INTO authors (name)
		VALUES ($1)
		RETURNING id
	`, newAuthor)
	
	var id int64

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, err
}

func (p *PGPool) DeleteAuthor(author string) (error) {

	_, err := p.pool.Exec(context.Background(), `
		DELETE FROM authors
		WHERE name = $1
	`, author)

	return err

}

