package domain

/*
type OrderStatus

- representar o status do pedido.

Status:
- PENDING
- PAID
- CANCELED
*/

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

func (s OrderStatus) String() string {
	return string(s)
}
