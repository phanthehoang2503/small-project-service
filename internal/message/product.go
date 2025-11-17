package message

type ProductMessage struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}
