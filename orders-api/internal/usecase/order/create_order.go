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
	"context"
	"time"

	"github.com/google/uuid"

	"orders-api/internal/domain"
	"orders-api/internal/messaging"
	"orders-api/internal/repository"
)

type CreateOrderUseCase struct {
	transactionManager repository.TransactionManager
	repositoryFactory  repository.RepositoryFactory
	eventPublisher     messaging.EventPublisher
}

func NewCreateOrderUseCase(
	transactionManager repository.TransactionManager,
	repositoryFactory repository.RepositoryFactory,
	publishers ...messaging.EventPublisher,
) *CreateOrderUseCase {

	var eventPublisher messaging.EventPublisher

	if len(publishers) > 0 {
		eventPublisher = publishers[0]
	}

	return &CreateOrderUseCase{
		transactionManager: transactionManager,
		repositoryFactory:  repositoryFactory,
		eventPublisher:     eventPublisher,
	}
}

func (uc *CreateOrderUseCase) Execute(order *domain.Order) error {

	if err := order.Validate(); err != nil {
		return err
	}

	err := uc.transactionManager.WithinTransaction(func(tx repository.DBTX) error {

		productRepository := uc.repositoryFactory.Product(tx)
		customerRepository := uc.repositoryFactory.Customer(tx)
		orderRepository := uc.repositoryFactory.Order(tx)

		if _, err := customerRepository.FindByID(order.CustomerID); err != nil {
			return err
		}

		products := make(map[string]*domain.Product)
		seen := make(map[string]struct{})

		for _, item := range order.Items {

			if _, exists := seen[item.ProductID]; exists {
				return domain.ErrDuplicatedProduct
			}

			seen[item.ProductID] = struct{}{}
		}

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

	if err != nil {
		return err
	}

	// A publicação é feita somente depois do commit da transação.
	if uc.eventPublisher != nil {

		event := messaging.Event{
			MessageID:     uuid.NewString(),
			EventType:     "OrderCreated",
			CorrelationID: uuid.NewString(),
			SagaID:        uuid.NewString(),
			OrderID:       order.ID,
			OccurredAt:    time.Now().UTC(),
		}

		if err := uc.eventPublisher.Publish(
			context.Background(),
			event,
		); err != nil {
			return err
		}
	}

	return nil
}
