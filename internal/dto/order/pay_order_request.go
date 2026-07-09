package order

/*
struct PayOrderRequest

Responsabilidades:

- receber o pagamento de um pedido.

Campos:
- orderId
*/

type PayOrderRequest struct {
	OrderID string `json:"orderId"`
}
