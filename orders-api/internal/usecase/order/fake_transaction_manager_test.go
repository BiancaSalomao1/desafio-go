package order

import "orders-api/internal/repository"

type FakeTransactionManager struct{}

func (f *FakeTransactionManager) WithinTransaction(
	fn func(tx repository.DBTX) error,
) error {

	return fn(nil)
}
