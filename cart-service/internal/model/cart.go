package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ProductID uint  `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
	Subtotal  int64 `json:"subtotal"`
}
