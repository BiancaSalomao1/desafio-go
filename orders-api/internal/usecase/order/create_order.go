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
	messaging "orders-api/internal/messaging"
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

		customerRepository := uc.repositoryFactory.Customer(tx)
		orderRepository := uc.repositoryFactory.Order(tx)

		if _, err := customerRepository.FindByID(order.CustomerID); err != nil {
			return err
		}

		seen := make(map[string]struct{})

		for _, item := range order.Items {

			if _, exists := seen[item.ProductID]; exists {
				return domain.ErrDuplicatedProduct
			}

			seen[item.ProductID] = struct{}{}

			if err := item.Validate(); err != nil {
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

	// A publicação ocorre somente depois do commit.
	if uc.eventPublisher != nil {

		items := make([]messaging.StockItem, 0, len(order.Items))

		for _, item := range order.Items {
			items = append(items, messaging.StockItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		event := messaging.Event{
			MessageID:     uuid.NewString(),
			EventType:     "ReserveStock",
			CorrelationID: uuid.NewString(),
			SagaID:        uuid.NewString(),
			OrderID:       order.ID,
			OccurredAt:    time.Now().UTC(),
			Data: messaging.ReserveStockData{
				OrderID: order.ID,
				Items:   items,
			},
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
