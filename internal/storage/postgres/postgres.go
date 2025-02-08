package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGPool struct {
	mu   sync.Mutex
	pool *pgxpool.Pool
}

func New(connStr string) (*PGPool, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PGPool{mu: sync.Mutex{}, pool: pool}, nil
}
