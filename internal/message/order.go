package message

type OrderMessage struct {
	ID        uint  `json:"id"`
	UserID    uint  `json:"user_id"`
	TotalCost int64 `json:"total_cost"`
}
