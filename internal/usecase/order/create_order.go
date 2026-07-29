package order

/*
struct CreateOrderUseCase
- criar pedido;
- validar cliente;
- validar produtos;
- validar estoque;
- reduzir estoque;
- salvar pedido.

Métodos:
- NewCreateOrderUseCase()
- Execute()
*/

import (
	"desafio-go/internal/domain"
	"desafio-go/internal/repository"
)

type CreateOrderUseCase struct {
	transactionManager repository.TransactionManager
	repositoryFactory  repository.RepositoryFactory
}

func NewCreateOrderUseCase(
	transactionManager repository.TransactionManager,
	repositoryFactory repository.RepositoryFactory,
) *CreateOrderUseCase {

	return &CreateOrderUseCase{
		transactionManager: transactionManager,
		repositoryFactory:  repositoryFactory,
	}
}

func (uc *CreateOrderUseCase) Execute(order *domain.Order) error {

	if err := order.Validate(); err != nil {
		return err
	}

	return uc.transactionManager.WithinTransaction(func(tx repository.DBTX) error {

		productRepository := uc.repositoryFactory.Product(tx)
		customerRepository := uc.repositoryFactory.Customer(tx)
		orderRepository := uc.repositoryFactory.Order(tx)

		if _, err := customerRepository.FindByID(order.CustomerID); err != nil {
			return err
		}

		products := make(map[string]*domain.Product)

		for i := range order.Items {

			product, err := productRepository.FindByID(order.Items[i].ProductID)
			if err != nil {
				return err
			}

			if err := product.ReduceStock(order.Items[i].Quantity); err != nil {
				return err
			}

			order.Items[i].Name = product.Name
			order.Items[i].Price = product.Price

			products[product.ID] = product
		}

		for _, product := range products {

			if err := productRepository.Update(product); err != nil {
				return err
			}
		}

		if err := orderRepository.Save(order); err != nil {
			return err
		}

		return nil
	})
}
