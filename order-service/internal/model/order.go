package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId uint   `json:"user_id"`
	Total  int64  `json:"total"`
	Status string `json:"status"`
	Item   []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderId   uint  `json:"order_id"`
	ProductId uint  `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
	Subtotal  int64 `json:"subtotal"`
}
