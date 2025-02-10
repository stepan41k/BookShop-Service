package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/stepan41k/testMidlware/internal/domain"
)

type event struct {
	ID      int    `db:"id"`
	Type    string `db:"event_type"`
	Payload string `db:"payload"`
}

func (p *PGPool) SaveEvent(tx pgx.Tx, eventType string, payload string) error {
	const op = "storage.postgres.genres.saveEvent"

	stmt, err := tx.Prepare(context.Background(), "my-query", "INSERT INTO events(event_type, payload) VALUES($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(context.Background(), stmt.SQL, eventType, payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PGPool) GetNewEvent(ctx context.Context) (domain.Event, error) {
	const op = "storage.posgres.genres.GetNewEvent"

	row := p.pool.QueryRow(context.Background(), `
		SELECT id, event_type, payload
		FROM events
		WHERE status = 'new'
		LIMIT 1	
	`)

	var evt event

	err := row.Scan(&evt.ID, &evt.Type, &evt.Payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Event{}, nil
		}

		return domain.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return domain.Event{
		ID:      evt.ID,
		Type:    evt.Type,
		Payload: evt.Payload,
	}, nil
}

func (p *PGPool) SetDone(id int) error {
	const op = "storage.postgres.genres.SetDone"

	_, err := p.pool.Exec(context.Background(), `
		UPDATE events SET status = 'done'
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
