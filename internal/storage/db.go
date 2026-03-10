package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) Pool() *pgxpool.Pool {
	return s.pool
}
