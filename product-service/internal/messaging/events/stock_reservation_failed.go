package events

type StockReservationFailedData struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}
