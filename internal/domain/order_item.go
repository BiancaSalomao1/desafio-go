package domain

/*
struct OrderItem

- identificar o item do pedido;
- armazenar produto;
- armazenar quantidade;
- armazenar preço da compra.

Métodos:
- construtor NewOrderItem()
- Subtotal()
- Validate()
*/

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

func NewOrderItem(productID, name string, price float64, quantity int) *OrderItem {
	return &OrderItem{
		ProductID: productID,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
	}
}

func (i *OrderItem) Validate() error {

	if i.ProductID == "" {
		return ErrProductInvalid
	}

	if i.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	return nil
}

func (i *OrderItem) Subtotal() float64 {
	return i.Price * float64(i.Quantity)
}
