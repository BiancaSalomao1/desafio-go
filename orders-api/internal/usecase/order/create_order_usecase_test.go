package order

import (
	"context"
	"errors"
	"testing"

	"orders-api/internal/domain"
	"orders-api/internal/messaging"
)

type eventPublisherSpy struct {
	calls  int
	events []messaging.Event
	err    error
}

func (p *eventPublisherSpy) Publish(
	_ context.Context,
	event messaging.Event,
) error {
	p.calls++
	p.events = append(p.events, event)

	return p.err
}

func TestCreateOrderUseCase_Execute(t *testing.T) {

	t.Run("should create order successfully", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		order := newOrder(
			customer.ID,
			product.ID,
			2,
		)

		factory := newFactory(product, customer)
		useCase := newUseCase(factory)

		err := useCase.Execute(order)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		// O orders-api não controla mais estoque.
		if product.Stock != 10 {
			t.Fatalf(
				"expected stock to remain 10, got %d",
				product.Stock,
			)
		}

		if factory.OrderRepository.saveCalls != 1 {
			t.Fatal("expected Save() to be called once")
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer FindByID() to be called once",
			)
		}

		if factory.ProductRepository.findCalls != 0 {
			t.Fatal(
				"product repository should not be called",
			)
		}

		if factory.ProductRepository.updateCalls != 0 {
			t.Fatal(
				"product repository should not be updated",
			)
		}
	})

	t.Run("should publish OrderCreated event", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		order := newOrder(
			customer.ID,
			product.ID,
			2,
		)

		factory := newFactory(product, customer)

		publisher := &eventPublisherSpy{}

		useCase := NewCreateOrderUseCase(
			&FakeTransactionManager{},
			factory,
			publisher,
		)

		err := useCase.Execute(order)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if publisher.calls != 1 {
			t.Fatalf(
				"expected Publish() to be called once, got %d",
				publisher.calls,
			)
		}

		event := publisher.events[0]

		if event.EventType != "ReserveStock" {
			t.Fatalf(
				"expected ReserveStock, got %s",
				event.EventType,
			)
		}

		if event.OrderID != order.ID {
			t.Fatalf(
				"expected order ID %s, got %s",
				order.ID,
				event.OrderID,
			)
		}

		if event.MessageID == "" {
			t.Fatal("expected MessageID")
		}

		if event.CorrelationID == "" {
			t.Fatal("expected CorrelationID")
		}

		if event.SagaID == "" {
			t.Fatal("expected SagaID")
		}

		data, ok := event.Data.(messaging.ReserveStockData)
		if !ok {
			t.Fatalf(
				"expected messaging.ReserveStockData, got %T",
				event.Data,
			)
		}

		if data.OrderID != order.ID {
			t.Fatalf(
				"expected event order ID %s, got %s",
				order.ID,
				data.OrderID,
			)
		}

		if len(data.Items) != 1 {
			t.Fatalf(
				"expected 1 item, got %d",
				len(data.Items),
			)
		}

		if data.Items[0].ProductID != product.ID {
			t.Fatalf(
				"expected product ID %s, got %s",
				product.ID,
				data.Items[0].ProductID,
			)
		}

		if data.Items[0].Quantity != 2 {
			t.Fatalf(
				"expected quantity 2, got %d",
				data.Items[0].Quantity,
			)
		}
	})

	t.Run("should return validation error", func(t *testing.T) {

		useCase := newUseCase(
			&FakeRepositoryFactory{},
		)

		order := &domain.Order{}

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("should return customer error", func(t *testing.T) {

		product := newProduct()

		factory := newFactory(product, nil)

		factory.CustomerRepository.findError =
			errors.New("customer not found")

		useCase := newUseCase(factory)

		order := newOrder(
			"customer",
			product.ID,
			1,
		)

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected customer error")
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer repository to be called once",
			)
		}

		if factory.ProductRepository.findCalls != 0 {
			t.Fatal(
				"product repository should not have been called",
			)
		}

		if factory.OrderRepository.saveCalls != 0 {
			t.Fatal(
				"order should not have been saved",
			)
		}
	})

	t.Run("should return save error", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		factory := newFactory(product, customer)

		factory.OrderRepository.saveError =
			errors.New("save error")

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			product.ID,
			2,
		)

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected save error")
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer repository to be called once",
			)
		}

		if factory.OrderRepository.saveCalls != 1 {
			t.Fatal(
				"expected order save attempt",
			)
		}

		if factory.ProductRepository.findCalls != 0 {
			t.Fatal(
				"product repository should not have been called",
			)
		}

		if factory.ProductRepository.updateCalls != 0 {
			t.Fatal(
				"product should not have been updated",
			)
		}
	})

	t.Run("should return publisher error", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		factory := newFactory(product, customer)

		publisher := &eventPublisherSpy{
			err: errors.New("publisher error"),
		}

		useCase := NewCreateOrderUseCase(
			&FakeTransactionManager{},
			factory,
			publisher,
		)

		order := newOrder(
			customer.ID,
			product.ID,
			1,
		)

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected publisher error")
		}

		if publisher.calls != 1 {
			t.Fatal(
				"expected Publish() to be called once",
			)
		}

		if factory.OrderRepository.saveCalls != 1 {
			t.Fatal(
				"expected order to be saved before publishing",
			)
		}
	})

	t.Run("should execute repositories in correct order", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		spy := &OrderSpy{}

		factory := newFactory(product, customer)

		factory.CustomerRepository.spy = spy
		factory.OrderRepository.spy = spy

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			product.ID,
			1,
		)

		err := useCase.Execute(order)

		if err != nil {
			t.Fatalf(
				"expected nil error, got %v",
				err,
			)
		}

		expected := []string{
			"customer.find",
			"order.save",
		}

		if !spy.Equals(expected) {
			t.Fatalf(
				"expected %v, got %v",
				expected,
				spy.Calls,
			)
		}
	})

	t.Run("should not interact with product repository", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		factory := newFactory(product, customer)

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			product.ID,
			1,
		)

		err := useCase.Execute(order)

		if err != nil {
			t.Fatalf(
				"expected nil error, got %v",
				err,
			)
		}

		if factory.ProductRepository.findCalls != 0 {
			t.Fatalf(
				"expected no product FindByID calls, got %d",
				factory.ProductRepository.findCalls,
			)
		}

		if factory.ProductRepository.updateCalls != 0 {
			t.Fatalf(
				"expected no product Update calls, got %d",
				factory.ProductRepository.updateCalls,
			)
		}
	})
}
