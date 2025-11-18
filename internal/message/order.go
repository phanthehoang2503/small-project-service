package message

type OrderRequested struct {
	UserID uint        `json:"user_id"`
	Items  []OrderItem `json:"items"`
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
