package database

/*
struct Database

- gerenciar conexão com PostgreSQL.

Métodos:
- New()
- Close()
*/

import (
	"context"
	"time"

	"orders-api/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
