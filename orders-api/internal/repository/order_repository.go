package repository

import "desafio-go/orders-api/internal/domain"

/*
interface OrderRepository

Responsabilidades:

- salvar pedido;
- atualizar pedido;
- remover pedido;
- buscar pedido por ID;
- listar pedidos.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

type OrderRepository interface {
	Save(order *domain.Order) error
	Update(order *domain.Order) error
	Delete(id string) error
	FindByID(id string) (*domain.Order, error)
	FindAll() ([]*domain.Order, error)
	List(limit, offset int) ([]*domain.Order, error)
}
