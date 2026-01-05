package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model      `swaggerignore:"true"`
	UUID            string      `json:"uuid" gorm:"size:36;uniqueIndex"`
	UserID          uint        `json:"user_id" gorm:"index;not null" example:"1"`
	Total           int64       `json:"total" example:"50000"`
	Status          string      `json:"status" example:"Pending"`
	ShippingAddress string      `json:"shipping_address" example:"123 Main St"`
	Items           []OrderItem `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderItem struct {
	gorm.Model `swaggerignore:"true"`
	OrderID    uint  `json:"order_id" gorm:"index;not null" example:"100"` // foreign key
	ProductID  uint  `json:"product_id" gorm:"index" example:"5"`
	Quantity   int   `json:"quantity" example:"2"`
	Price      int64 `json:"price" example:"25000"`
	Subtotal   int64 `json:"subtotal" example:"50000"`
}
