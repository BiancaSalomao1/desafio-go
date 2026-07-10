package repository

import "github.com/jackc/pgx/v5"

/*
TransactionManager

Responsabilidades:
- executar operações dentro de uma transação.

Métodos:
- WithinTransaction()
*/

type TransactionManager interface {
	WithinTransaction(func(tx pgx.Tx) error) error
}
