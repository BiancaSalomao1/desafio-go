package events

type ReleaseStockData struct {
	OrderID string      `json:"order_id"`
	Items   []StockItem `json:"items"`
}
