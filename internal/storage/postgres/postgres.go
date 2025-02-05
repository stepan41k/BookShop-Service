package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGPool struct {
	mu sync.Mutex
	pool *pgxpool.Pool
}

func New(connStr string) (*PGPool, error) {
	Conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &PGPool{mu: sync.Mutex{}, pool: Conn}, nil
}