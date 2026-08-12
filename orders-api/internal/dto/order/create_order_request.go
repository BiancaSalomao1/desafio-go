package order

/*
struct CreateOrderRequest

Responsabilidades:
- receber os dados para criação de um pedido.

Campos:
- customerId
- items

Métodos:
- nenhum
*/

type CreateOrderRequest struct {
	CustomerID string                   `json:"customerId"`
	Items      []CreateOrderItemRequest `json:"items"`
}
