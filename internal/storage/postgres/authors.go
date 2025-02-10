package postgres

import (
	"context"
	"fmt"
)

const(
	statusAuthorCreated = "AuthorCreated"
	statusAuthorDeleted = "AuthorDeleted"
)

func (p *PGPool) SaveAuthor(newAuthor string) (id int64, err error) {
	const op = "storage.postgres.authors.SaveAuthor"

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func()  {
		if err != nil{
			_ = tx.Rollback(context.Background())
			return
		}	

		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	err = p.pool.QueryRow(context.Background(), `
		INSERT INTO authors (author)
		VALUES ($1)
		RETURNING id
	`).Scan(&id)
	
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	eventPayload := fmt.Sprintf(
		`{"id": %d, "author": %s}`,
		id,
		newAuthor,
	)

	if err := p.SaveEvent(tx, statusAuthorCreated, eventPayload); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, err
}

func (p *PGPool) DeleteAuthor(author string) (err error) {
	const op = "storage.postgres.authors.DeleteAuthor"

	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func()  {
		if err != nil {
			_ = tx.Rollback(context.Background())
			return	
		}

		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	_, err = p.pool.Exec(context.Background(), `
		DELETE FROM authors
		WHERE name = $1
	`, author)
	
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	eventPayload := fmt.Sprintf(
		`{"author": %s}`,
		author,
	)

	if err := p.SaveEvent(tx, statusAuthorDeleted, eventPayload); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

