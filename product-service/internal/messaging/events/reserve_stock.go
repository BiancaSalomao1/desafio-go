package events

type StockItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveStockData struct {
	OrderID string      `json:"order_id"`
	Items   []StockItem `json:"items"`
}
