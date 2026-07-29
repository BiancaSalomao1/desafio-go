package repository

/*
TransactionManager

Responsabilidades:
- executar operações dentro de uma transação.
*/

type TransactionManager interface {
	WithinTransaction(func(tx DBTX) error) error
}
