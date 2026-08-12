package order

/*
Helpers utilizados pelos testes do módulo Order.

Responsabilidades:
- criar entidades válidas;
- criar FakeRepositoryFactory;
- criar CreateOrderUseCase.
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

func newProduct() *domain.Product {
	return &domain.Product{
		ID:    "p1",
		Name:  "Notebook",
		Price: 5000,
		Stock: 10,
	}
}

func newCustomer() *domain.Customer {
	return &domain.Customer{
		ID:    "c1",
		Name:  "João da Silva",
		Email: "joao@email.com",
	}
}

func newOrder(customerID, productID string, quantity int) *domain.Order {

	order := domain.NewOrder(
		"o1",
		customerID,
	)

	_ = order.AddItem(domain.OrderItem{
		ID:        "i1",
		ProductID: productID,
		Name:      "Notebook",
		Price:     5000,
		Quantity:  quantity,
	})

	return order
}

func newFactory(
	product *domain.Product,
	customer *domain.Customer,
) *FakeRepositoryFactory {

	return &FakeRepositoryFactory{
		ProductRepository: &FakeProductRepository{
			product: product,
		},
		CustomerRepository: &FakeCustomerRepository{
			customer: customer,
		},
		OrderRepository: &FakeOrderRepository{},
	}
}

func newUseCase(
	factory repository.RepositoryFactory,
) *CreateOrderUseCase {

	return NewCreateOrderUseCase(
		&FakeTransactionManager{},
		factory,
	)
}
