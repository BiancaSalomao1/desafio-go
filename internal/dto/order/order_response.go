package order

/*
struct OrderResponse

Responsabilidades:

- retornar os dados de um pedido.

Campos:
- id
- customerId
- items
- status
- total
*/

type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customerId"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	Total      float64             `json:"total"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}
