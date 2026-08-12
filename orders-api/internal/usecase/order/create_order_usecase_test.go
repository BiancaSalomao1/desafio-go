package order

/*
Testes do CreateOrderUseCase.

Responsabilidades:
- validar criação de pedidos;
- validar regras de negócio;
- validar interação entre repositórios.

Cenários:
- pedido criado com sucesso;
- pedido inválido;
- cliente inexistente;
- produto inexistente.
*/

import (
	"errors"
	"testing"

	"desafio-go/orders-api/internal/domain"
)

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

		if product.Stock != 8 {
			t.Fatalf(
				"expected stock 8, got %d",
				product.Stock,
			)
		}

		if factory.OrderRepository.saveCalls != 1 {
			t.Fatal(
				"expected Save() to be called once",
			)
		}

		if factory.ProductRepository.updateCalls != 1 {
			t.Fatal(
				"expected Update() to be called once",
			)
		}

		if factory.ProductRepository.findCalls != 1 {
			t.Fatal(
				"expected FindByID() to be called once",
			)
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer FindByID() to be called once",
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
	t.Run("should return product error", func(t *testing.T) {

		customer := newCustomer()

		factory := newFactory(nil, customer)

		factory.ProductRepository.findError =
			errors.New("product not found")

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			"product-1",
			1,
		)

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected product error")
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer repository to be called once",
			)
		}

		if factory.ProductRepository.findCalls != 1 {
			t.Fatal(
				"expected product repository to be called once",
			)
		}

		if factory.ProductRepository.updateCalls != 0 {
			t.Fatal(
				"product should not have been updated",
			)
		}

		if factory.OrderRepository.saveCalls != 0 {
			t.Fatal(
				"order should not have been saved",
			)
		}
	})

	t.Run("should return update error", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		factory := newFactory(product, customer)

		factory.ProductRepository.updateError =
			errors.New("update error")

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			product.ID,
			2,
		)

		err := useCase.Execute(order)

		if err == nil {
			t.Fatal("expected update error")
		}

		if factory.CustomerRepository.findCalls != 1 {
			t.Fatal(
				"expected customer repository to be called once",
			)
		}

		if factory.ProductRepository.findCalls != 1 {
			t.Fatal(
				"expected product repository to be called once",
			)
		}

		if factory.ProductRepository.updateCalls != 1 {
			t.Fatal(
				"expected Update() to be called once",
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

		if factory.ProductRepository.findCalls != 1 {
			t.Fatal(
				"expected product repository to be called once",
			)
		}

		if factory.ProductRepository.updateCalls != 1 {
			t.Fatal(
				"expected product update",
			)
		}

		if factory.OrderRepository.saveCalls != 1 {
			t.Fatal(
				"expected order save attempt",
			)
		}
	})
	t.Run("should execute repositories in correct order", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		spy := &OrderSpy{}

		factory := newFactory(product, customer)

		factory.ProductRepository.spy = spy
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
			t.Fatalf("expected nil error, got %v", err)
		}

		expected := []string{
			"customer.find",
			"product.find",
			"product.update",
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

	t.Run("should satisfy mock expectations", func(t *testing.T) {

		product := newProduct()
		customer := newCustomer()

		mock := &OrderMock{
			ExpectedCustomerFind:  1,
			ExpectedProductFind:   1,
			ExpectedProductUpdate: 1,
			ExpectedOrderSave:     1,
		}

		factory := newFactory(product, customer)

		factory.ProductRepository.mock = mock
		factory.CustomerRepository.mock = mock
		factory.OrderRepository.mock = mock

		useCase := newUseCase(factory)

		order := newOrder(
			customer.ID,
			product.ID,
			1,
		)

		err := useCase.Execute(order)

		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if err := mock.Verify(); err != nil {
			t.Fatal(err)
		}
	})

}
