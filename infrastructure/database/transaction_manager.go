package database

/*
struct TransactionManager

Responsabilidades:
- executar operações dentro de uma transação.

Campos:
- pool

Métodos:
- NewTransactionManager()
- WithinTransaction()
*/

import (
	"context"

	"desafio-go/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(
	pool *pgxpool.Pool,
) repository.TransactionManager {

	return &TransactionManager{
		pool: pool,
	}
}

func (tm *TransactionManager) WithinTransaction(
	fn func(tx pgx.Tx) error,
) error {

	ctx := context.Background()

	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
