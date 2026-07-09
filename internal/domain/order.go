package domain

/*
struct Order

- identificar o pedido;
- armazenar cliente;
- armazenar itens;
- armazenar status.

Métodos:
- construtor NewOrder()
- AddItem()
- RemoveItem()
- Total()
- Pay()
- Cancel()
- Validate()
*/

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     OrderStatus `json:"status"`
}

// NewOrder cria um novo pedido com status PENDING.
func NewOrder(id string, customerID string) *Order {
	return &Order{
		ID:         id,
		CustomerID: customerID,
		Status:     OrderStatusPending,
		Items:      []OrderItem{},
	}
}

func (o *Order) AddItem(item OrderItem) error {
	if item.Validate() != nil {
		return item.Validate()
	}

	o.Items = append(o.Items, item)

	return nil
}

func (o *Order) RemoveItem(productID string) error {

	for i, item := range o.Items {

		if item.ProductID == productID {

			o.Items = append(o.Items[:i], o.Items[i+1:]...)

			return nil
		}
	}

	return ErrProductNotFound
}

func (o *Order) Total() float64 {

	total := 0.0

	for _, item := range o.Items {
		total += item.Subtotal()
	}

	return total
}

func (o *Order) Pay() error {

	if o.Status != OrderStatusPending {
		return ErrOrderStatusInvalid
	}

	o.Status = OrderStatusPaid

	return nil
}

func (o *Order) Cancel() error {

	if o.Status != OrderStatusPending {
		return ErrOrderStatusInvalid
	}

	o.Status = OrderStatusCanceled

	return nil
}

func (o *Order) Validate() error {

	if o.CustomerID == "" {
		return ErrCustomerInvalid
	}

	if len(o.Items) == 0 {
		return ErrEmptyOrder
	}

	return nil
}
