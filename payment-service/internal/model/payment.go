package model

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	OrderUUID string `json:"order_uuid" gorm:"size:36;index;uniqueIndex"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency" gorm:"size:8"`
	Provider  string `json:"provider" gorm:"size:32"` // "mock"
	Status    string `json:"status" gorm:"size:32"`   // PENDING, SUCCEEDED, FAILED
}
