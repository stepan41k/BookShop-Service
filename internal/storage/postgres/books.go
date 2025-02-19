package postgres

import (
	"context"
	"fmt"

	"github.com/stepan41k/testMidlware/internal/domain"
)

const (
	statusBookCreated = "BookCreated"
	statusBookDeleted = "BookDeleted"
	statusBookUpdated = "BookUpdated"
)

func (p *PGPool) GetBooksByAuthor(author string) ([]domain.Book, error) {
	
	rows, err := p.pool.Query(context.Background(), `
		SELECT id, name, author_id, genre_id, price
		FROM books
		WHERE author = $1;
	`, author)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []domain.Book
	for rows.Next() {
		var item domain.Book
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Author,
			&item.Genre,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return data, nil
}

func (p *PGPool) GetBooksByGenre(genre string) ([]domain.Book, error) {
	const op = "storage.postgres.books.GetBooksByGenre"

	rows, err := p.pool.Query(context.Background(), `
		SELECT id, name, author_id, genre_id, price
		FROM books
		WHERE genre = $1;
	`, genre)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var data []domain.Book
	for rows.Next() {
		var item domain.Book
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Author,
			&item.Genre,
			&item.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		data = append(data, item)
	}

	return data, nil
}

func (p *PGPool) GetBookByName(name string) (book domain.Book, err error) {
	const op = "storage.postgres.books.GetBookByName"

	err = p.pool.QueryRow(context.Background(), `
		SELECT id, name, author_id, genre_id, price
		FROM books
		WHERE name = $1;
	`,	name).Scan(
		&book.ID,
		&book.Name,
		&book.Author,
		&book.Genre,
		&book.Price,
	)

	if err != nil {
		return domain.Book{}, fmt.Errorf("%s: %w", op, err)
	}

	return book, nil
}

func (p *PGPool) SaveBook(item domain.Book) (id int64, err error) {
	const op = "storage.postgres.books.SaveBook"

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

	err = p.pool.QueryRow(context.Background(), `
		INSERT INTO books (name, author_id, genre_id, price)
		VALUES (
			$1,
			(SELECT id FROM authors WHERE author = $2),
			(SELECT id FROM genres WHERE genre = $3),
			$4
		)
		RETURNING id;`,
		item.Name,
		item.Author,
		item.Genre,
		item.Price,
		).Scan(&id)
	
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	eventPayload := fmt.Sprintf(
		`{"id:" %d, "name:", %s, "author:", %s, "genre:", %s, "price:", %s}`,
		id,
		item.Name,
		item.Author,
		item.Genre,
		item.Price,
	)

	if err := p.SaveEvent(tx, statusBookCreated, eventPayload); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	
	return id, nil
}

func (p * PGPool) DeleteBook(name string) (error) {
	const op = "storage.postgres.books.DeleteBook"

	_, err := p.pool.Exec(context.Background(), `
		DELETE FROM books
		WHERE name = $1
		`, name)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PGPool) UpdateBook(oldBook string, newBook domain.Book)(error) {
	const op = "storage.postgres.books.UpdateBook"

	_, err := p.pool.Exec(context.Background(), `
	UPDATE books
	SET name = $1, author_id = $2, genre_id = $3, price = $4
	WHERE name = $5`,
	newBook.Name, newBook.Author, newBook.Genre, newBook.Price, oldBook)
	
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
		
	return nil
}