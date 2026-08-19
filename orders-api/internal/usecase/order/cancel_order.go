package order

import (
	"context"

	"orders-api/internal/messaging"
	"orders-api/internal/repository"
)

type CancelOrderUseCase struct {
	transactionManager repository.TransactionManager
	repositoryFactory  repository.RepositoryFactory
	publishers         []messaging.EventPublisher
}

func NewCancelOrderUseCase(
	transactionManager repository.TransactionManager,
	repositoryFactory repository.RepositoryFactory,
	publishers ...messaging.EventPublisher,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		transactionManager: transactionManager,
		repositoryFactory:  repositoryFactory,
		publishers:         publishers,
	}
}

func (uc *CancelOrderUseCase) Execute(id string) error {
	var event messaging.Event

	err := uc.transactionManager.WithinTransaction(func(tx repository.DBTX) error {
		orderRepository := uc.repositoryFactory.Order(tx)

		order, err := orderRepository.FindByID(id)
		if err != nil {
			return err
		}

		if err := order.Cancel(); err != nil {
			return err
		}

		if err := orderRepository.Update(order); err != nil {
			return err
		}

		items := make([]messaging.StockItem, 0, len(order.Items))

		for _, item := range order.Items {
			items = append(items, messaging.StockItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		event = messaging.Event{
			EventType: "ReleaseStock",
			OrderID:   order.ID,
			Data: struct {
				Items []messaging.StockItem `json:"items"`
			}{
				Items: items,
			},
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, publisher := range uc.publishers {
		if err := publisher.Publish(context.Background(), event); err != nil {
			return err
		}
	}

	return nil
}
