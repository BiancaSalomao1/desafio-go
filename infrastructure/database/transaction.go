package database

/*
TransactionManager

Responsabilidades:
- iniciar transações;
- confirmar transações;
- desfazer transações.

Métodos:
- WithinTransaction()
*/

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (db *Database) WithinTransaction(
	fn func(tx pgx.Tx) error,
) error {

	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
