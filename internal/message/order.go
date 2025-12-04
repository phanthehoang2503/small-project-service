package message

type OrderRequested struct {
	CorrelationID string      `json:"correlation_id"`
	OrderUUID     string      `json:"order_uuid"`
	UserID        uint        `json:"user_id"`
	Total         int64       `json:"total"`
	Currency      string      `json:"currency"`
	Items         []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderCreated struct {
	OrderID uint  `json:"order_id"`
	UserID  uint  `json:"user_id"`
	Total   int64 `json:"total"`
}

type StockFailed struct {
	OrderUUID string `json:"order_uuid"`
	Reason    string `json:"reason"`
}
