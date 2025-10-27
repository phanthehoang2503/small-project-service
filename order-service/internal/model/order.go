package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId uint   `json:"user_id"`
	Total  int64  `json:"total"`
	Status string `json:"status"`
	Item[]
}
