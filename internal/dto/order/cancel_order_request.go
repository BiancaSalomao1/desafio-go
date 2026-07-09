package order

/*
struct CancelOrderRequest

Responsabilidades:

- receber o cancelamento de um pedido.

Campos:
- orderId
*/

type CancelOrderRequest struct {
	OrderID string `json:"orderId"`
}
