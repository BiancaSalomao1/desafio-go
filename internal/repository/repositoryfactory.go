package repository

/*
RepositoryFactory

Responsabilidades:
- criar repositories utilizando um DBTX.
*/

type RepositoryFactory interface {
	Product(db DBTX) ProductRepository
	Customer(db DBTX) CustomerRepository
	User(db DBTX) UserRepository
	Order(db DBTX) OrderRepository
}
