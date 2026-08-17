package testutil

import (
	"orders-api/internal/domain"
)

/******************************
* Product
******************************/

type CreateProductMock func(*domain.Product) error

func (f CreateProductMock) Execute(product *domain.Product) error {
	return f(product)
}

type GetProductMock func(string) (*domain.Product, error)

func (f GetProductMock) Execute(id string) (*domain.Product, error) {
	return f(id)
}

type ListProductsMock func() ([]*domain.Product, error)

func (f ListProductsMock) Execute() ([]*domain.Product, error) {
	return f()
}

type UpdateProductMock func(*domain.Product) error

func (f UpdateProductMock) Execute(product *domain.Product) error {
	return f(product)
}

type DeleteProductMock func(string) error

func (f DeleteProductMock) Execute(id string) error {
	return f(id)
}

/******************************
* Customer
******************************/

type CreateCustomerMock func(*domain.Customer) error

func (f CreateCustomerMock) Execute(customer *domain.Customer) error {
	return f(customer)
}

type GetCustomerMock func(string) (*domain.Customer, error)

func (f GetCustomerMock) Execute(id string) (*domain.Customer, error) {
	return f(id)
}

type ListCustomersMock func() ([]*domain.Customer, error)

func (f ListCustomersMock) Execute() ([]*domain.Customer, error) {
	return f()
}

type UpdateCustomerMock func(*domain.Customer) error

func (f UpdateCustomerMock) Execute(customer *domain.Customer) error {
	return f(customer)
}

type DeleteCustomerMock func(string) error

func (f DeleteCustomerMock) Execute(id string) error {
	return f(id)
}

/******************************
* User
******************************/

type CreateUserMock func(*domain.User) error

func (f CreateUserMock) Execute(user *domain.User) error {
	return f(user)
}

type GetUserMock func(string) (*domain.User, error)

func (f GetUserMock) Execute(id string) (*domain.User, error) {
	return f(id)
}

type ListUsersMock func() ([]*domain.User, error)

func (f ListUsersMock) Execute() ([]*domain.User, error) {
	return f()
}

type UpdateUserMock func(*domain.User) error

func (f UpdateUserMock) Execute(user *domain.User) error {
	return f(user)
}

type DeleteUserMock func(string) error

func (f DeleteUserMock) Execute(id string) error {
	return f(id)
}

/******************************
* Order
******************************/

type CreateOrderMock func(*domain.Order) error

func (f CreateOrderMock) Execute(order *domain.Order) error {
	return f(order)
}

type GetOrderMock func(string) (*domain.Order, error)

func (f GetOrderMock) Execute(id string) (*domain.Order, error) {
	return f(id)
}

type ListOrdersMock func(int, int) ([]*domain.Order, error)

func (f ListOrdersMock) Execute(page, limit int) ([]*domain.Order, error) {
	return f(page, limit)
}

type PayOrderMock func(string) error

func (f PayOrderMock) Execute(id string) error {
	return f(id)
}

type CancelOrderMock func(string) error

func (f CancelOrderMock) Execute(id string) error {
	return f(id)
}

/******************************
* Auth
******************************/

type LoginMock func(string, string) (string, error)

func (f LoginMock) Execute(email, password string) (string, error) {
	return f(email, password)
}
