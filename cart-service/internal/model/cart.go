package model

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model `swaggerignore:"true"`
	UserID     uint  `json:"user_id" gorm:"index"`
	ProductID  uint  `json:"product_id"`
	Quantity   int   `json:"quantity"`
	Price      int64 `json:"price"`
	Subtotal   int64 `json:"subtotal"`
}
