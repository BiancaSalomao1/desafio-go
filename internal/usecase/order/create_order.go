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
	orderRepository    repository.OrderRepository
	productRepository  repository.ProductRepository
	customerRepository repository.CustomerRepository
}

func NewCreateOrderUseCase(
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
	customerRepository repository.CustomerRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepository:    orderRepository,
		productRepository:  productRepository,
		customerRepository: customerRepository,
	}
}

func (uc *CreateOrderUseCase) Execute(order *domain.Order) error {
	//validações
	//  pedido
	if err := order.Validate(); err != nil {
		return err
	}

	// cliente existe
	if _, err := uc.customerRepository.FindByID(order.CustomerID); err != nil {
		return err
	}

	// atualiza apenas uma vez
	products := make(map[string]*domain.Product)

	// valida produtos e reduz estoque
	for i := range order.Items {

		product, err := uc.productRepository.FindByID(order.Items[i].ProductID)
		if err != nil {
			return err
		}

		if err := product.ReduceStock(order.Items[i].Quantity); err != nil {
			return err
		}

		// **congela** os dados do produto no momento da compra
		order.Items[i].Name = product.Name
		order.Items[i].Price = product.Price

		products[product.ID] = product
	}

	// atualiza o estoque
	for _, product := range products {

		if err := uc.productRepository.Update(product); err != nil {
			return err
		}
	}

	// salva
	if err := uc.orderRepository.Save(order); err != nil {
		return err
	}

	return nil
}
