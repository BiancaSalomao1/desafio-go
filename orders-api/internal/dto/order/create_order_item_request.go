package order

/*
struct CreateOrderItemRequest

Responsabilidades:
- receber os dados de um item do pedido.

Campos:
- productId
- quantity

*/

type CreateOrderItemRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
