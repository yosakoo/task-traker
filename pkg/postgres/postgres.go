package postgres

import (
	"context"
	"fmt"
	"github.com/yosakoo/task-traker/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool  *pgxpool.Pool
	log logger.Interface
}

func New(dbURL string, log *logger.Logger) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return &Storage{
		Pool:  pool,
		log: log,
	}, nil
}

func (p *Storage) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}