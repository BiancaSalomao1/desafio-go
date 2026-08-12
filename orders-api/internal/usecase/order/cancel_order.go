package order

/*
struct CancelOrderUseCase

Responsabilidades:
- cancelar pedido;
- devolver estoque;
- atualizar pedido;
- executar toda a operação em uma transação.

Campos:
- transactionManager
- repositoryFactory

Métodos:
- NewCancelOrderUseCase()
- Execute()
*/

import (
	"desafio-go/orders-api/internal/repository"
)

type CancelOrderUseCase struct {
	transactionManager repository.TransactionManager
	repositoryFactory  repository.RepositoryFactory
}

func NewCancelOrderUseCase(
	transactionManager repository.TransactionManager,
	repositoryFactory repository.RepositoryFactory,
) *CancelOrderUseCase {

	return &CancelOrderUseCase{
		transactionManager: transactionManager,
		repositoryFactory:  repositoryFactory,
	}
}

func (uc *CancelOrderUseCase) Execute(id string) error {

	return uc.transactionManager.WithinTransaction(func(tx repository.DBTX) error {

		orderRepository := uc.repositoryFactory.Order(tx)
		productRepository := uc.repositoryFactory.Product(tx)

		order, err := orderRepository.FindByID(id)
		if err != nil {
			return err
		}

		if err := order.Cancel(); err != nil {
			return err
		}

		for _, item := range order.Items {

			product, err := productRepository.FindByID(item.ProductID)
			if err != nil {
				return err
			}

			product.IncreaseStock(item.Quantity)

			if err := productRepository.Update(product); err != nil {
				return err
			}
		}

		if err := orderRepository.Update(order); err != nil {
			return err
		}

		return nil
	})
}
