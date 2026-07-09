package order

/*
struct CreateOrderRequest

Responsabilidades:

- receber os dados para criação de um pedido.

Campos:
- customerId
- items
*/

type CreateOrderRequest struct {
	CustomerID string                   `json:"customerId"`
	Items      []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
