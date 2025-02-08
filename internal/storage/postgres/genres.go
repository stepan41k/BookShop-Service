package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *PGPool) SaveGenre(newGenre string) (int64, error) {

	const op = "storage.postgres.genres.SaveGenre"

	row := p.pool.QueryRow(context.Background(), `
	INSERT INTO genres (genre)
	VALUES($1)
	RETURNING id;
	`, newGenre)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}


	return id, err

}

func (p *PGPool) DeleteGenre(genre string) (error) {
	
	const op = "storage.postgres.genres.DeleteGenre"

	_, err := p.pool.Exec(context.Background(), `
		DELETE FROM genres
		WHERE genre = $1
	`, genre)

	if err != nil {
		
		//TODO:
		
		return fmt.Errorf("%s: %w", op, err)
	}

	return err
}

