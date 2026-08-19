package messaging

type StockItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderCreatedData struct {
	OrderID string      `json:"order_id"`
	Items   []StockItem `json:"items"`
}

type ReserveStockData struct {
	OrderID string      `json:"order_id"`
	Items   []StockItem `json:"items"`
}
