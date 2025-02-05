package repository

import (
	"context"
	"github.com/stepan41k/testMidlware/internal/pkg/models"
)

func (repo *PGRepo) CreateGenre(newGenre models.Genre) (id int ,err error) {

	err = repo.pool.QueryRow(context.Background(), `
	INSERT INTO genres (genre)
	VALUES($1)
	RETURNING id;
	`, newGenre.Genre).Scan(
		id,
	)
	return id, err

}

func (repo *PGRepo) ReadGenre() ([]models.Genre, error) {
	rows, err := repo.pool.Query(context.Background(), `
		SELECT id, genre
		FROM genres
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data []models.Genre

	for rows.Next() {
		var item models.Genre
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

func (repo *PGRepo) ReadGenreById(id int) (result models.Genre, err error) {
	err = repo.pool.QueryRow(context.Background(), `
		SELECT id, genre
		FROM genres
		WHERE id = $1
	`, id).Scan(
		&result.ID,
		&result.Genre,
	)

	return result, err
}

func (repo *PGRepo) UpdateGenre(genre string, newGenre models.Genre) (error) {
	_, err := repo.pool.Exec(context.Background(), `
		UPDATE genres
		SET genre = $1
		WHERE genre = $2
	`, newGenre.Genre, genre)

	return err
}

func (repo *PGRepo) DeleteGenre(genre string) (error) {
	
	_, err := repo.pool.Exec(context.Background(), `
		DELETE FROM genres
		WHERE genre = $1
	`, genre)

	return err
}

