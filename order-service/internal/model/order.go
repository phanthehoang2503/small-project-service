package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model      `swaggerignore:"true"`
	UserID          uint        `json:"user_id" gorm:"index;not null"`
	Total           int64       `json:"total"`
	Status          string      `json:"status"`
	ShippingAddress string      `json:"shipping_address"`
	Items           []OrderItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderItem struct {
	gorm.Model `swaggerignore:"true"`
	OrderID    uint  `json:"order_id" gorm:"index;not null"` // foreign key
	ProductID  uint  `json:"product_id" gorm:"index"`
	Quantity   int   `json:"quantity"`
	Price      int64 `json:"price"`
	Subtotal   int64 `json:"subtotal"`
}
