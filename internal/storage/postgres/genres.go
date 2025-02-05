package postgres

import (
	"context"
	"github.com/stepan41k/testMidlware/internal/storage"
)

func (p *PGPool) CreateGenre(newGenre storage.Genre) (id int ,err error) {

	err = p.pool.QueryRow(context.Background(), `
	INSERT INTO genres (genre)
	VALUES($1)
	RETURNING id;
	`, newGenre.Genre).Scan(
		id,
	)
	return id, err

}

func (p *PGPool) ReadGenre() ([]storage.Genre, error) {
	rows, err := p.pool.Query(context.Background(), `
		SELECT id, genre
		FROM genres
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []storage.Genre

	for rows.Next() {
		var item storage.Genre
		err = rows.Scan(
			&item.ID,
			&item.Genre,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, nil
}

func (p *PGPool) ReadGenreById(id int) (result storage.Genre, err error) {
	err = p.pool.QueryRow(context.Background(), `
		SELECT id, genre
		FROM genres
		WHERE id = $1
	`, id).Scan(
		&result.ID,
		&result.Genre,
	)

	return result, err
}

func (p *PGPool) UpdateGenre(genre string, newGenre storage.Genre) (error) {
	_, err := p.pool.Exec(context.Background(), `
		UPDATE genres
		SET genre = $1
		WHERE genre = $2
	`, newGenre.Genre, genre)

	return err
}

func (p *PGPool) DeleteGenre(genre string) (error) {
	
	_, err := p.pool.Exec(context.Background(), `
		DELETE FROM genres
		WHERE genre = $1
	`, genre)

	return err
}

