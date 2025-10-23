package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}
