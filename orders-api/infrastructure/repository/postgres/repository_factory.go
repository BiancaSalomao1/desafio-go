package postgres

import "orders-api/internal/repository"

type RepositoryFactory struct{}

func NewRepositoryFactory() repository.RepositoryFactory {
	return &RepositoryFactory{}
}

func (f *RepositoryFactory) Product(db repository.DBTX) repository.ProductRepository {
	return NewProductRepository(db)
}

func (f *RepositoryFactory) Customer(db repository.DBTX) repository.CustomerRepository {
	return NewCustomerRepository(db)
}

func (f *RepositoryFactory) User(db repository.DBTX) repository.UserRepository {
	return NewUserRepository(db)
}

func (f *RepositoryFactory) Order(db repository.DBTX) repository.OrderRepository {
	return NewOrderRepository(db)
}
