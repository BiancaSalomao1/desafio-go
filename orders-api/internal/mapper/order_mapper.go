package mapper

import (
	"orders-api/internal/domain"
	orderdto "orders-api/internal/dto/order"

	"github.com/google/uuid"
)

func ToOrder(
	request orderdto.CreateOrderRequest,
) *domain.Order {

	order := domain.NewOrder(
		uuid.NewString(),
		request.CustomerID,
	)

	for _, item := range request.Items {

		order.AddItem(
			*domain.NewOrderItem(
				item.ProductID,
				"",
				0,
				item.Quantity,
			),
		)
	}

	return order
}

func ToOrderResponse(
	order *domain.Order,
) orderdto.OrderResponse {

	items := make(
		[]orderdto.OrderItemResponse,
		0,
		len(order.Items),
	)

	for _, item := range order.Items {

		items = append(
			items,
			orderdto.OrderItemResponse{
				ID:        item.ID,
				ProductID: item.ProductID,
				Name:      item.Name,
				Price:     item.Price,
				Quantity:  item.Quantity,
				Subtotal:  item.Subtotal(),
			},
		)
	}

	return orderdto.OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     string(order.Status),
		Items:      items,
		Total:      order.Total(),
	}
}

func ToOrderResponseList(
	orders []*domain.Order,
) []orderdto.OrderResponse {

	response := make(
		[]orderdto.OrderResponse,
		0,
		len(orders),
	)

	for _, order := range orders {

		response = append(
			response,
			ToOrderResponse(order),
		)
	}

	return response
}
