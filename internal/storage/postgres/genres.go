package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	statusGenreCreated = "GenreCreated"
	statusGenreDeleted = "GenreDeleted"
)

func (p *PGPool) SaveGenre(newGenre string) (id int64, err error) {
	const op = "storage.postgres.genres.SaveGenre"

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}

		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	err = tx.QueryRow(context.Background(), `
	INSERT INTO genres (genre)
	VALUES($1)
	RETURNING id;
	`, newGenre).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	eventPayload := fmt.Sprintf(`
	{"id": %d, "genre": "%s"}`,
		id,
		newGenre)

	if err := p.SaveEvent(tx, statusGenreCreated, eventPayload); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, err

}

func (p *PGPool) DeleteGenre(genre string) error {
	const op = "storage.postgres.genres.DeleteGenre"

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	
	defer func ()  {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return
		}	

		comitErr := tx.Commit(context.Background())
		if comitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()
		_, err = p.pool.Exec(context.Background(), `
		DELETE FROM genres
		WHERE genre = $1
	`, genre)

	if err != nil {

		//TODO:

		return fmt.Errorf("%s: %w", op, err)
	}

	eventPayload := fmt.Sprintf(`{"genre": %s}`,
		genre,
	)

	if err := p.SaveEvent(tx, statusGenreDeleted, eventPayload); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
