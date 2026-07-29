package order

import (
	"desafio-go/internal/repository"
)

type FakeRepositoryFactory struct {
	ProductRepository  *FakeProductRepository
	CustomerRepository *FakeCustomerRepository
	OrderRepository    *FakeOrderRepository
}

func (f *FakeRepositoryFactory) Product(
	tx repository.DBTX,
) repository.ProductRepository {

	return f.ProductRepository
}

func (f *FakeRepositoryFactory) Customer(
	tx repository.DBTX,
) repository.CustomerRepository {

	return f.CustomerRepository
}

func (f *FakeRepositoryFactory) Order(
	tx repository.DBTX,
) repository.OrderRepository {

	return f.OrderRepository
}

func (f *FakeRepositoryFactory) User(
	tx repository.DBTX,
) repository.UserRepository {
	return nil
}
