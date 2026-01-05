package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model `swaggerignore:"true"`
	OrderUUID  string `json:"order_uuid" gorm:"size:36;index;uniqueIndex" example:"2133223-123213-3213"`
	Amount     int64  `json:"amount" example:"50000"`
	Currency   string `json:"currency" gorm:"size:8" example:"VND"`
	Provider   string `json:"provider" gorm:"size:32" example:"momo"`
	Status     string `json:"status" gorm:"size:32" example:"paid"` // PENDING, SUCCEEDED, FAILED
}
