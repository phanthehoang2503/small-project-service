package model

import "gorm.io/gorm"

type ProductSnapshot struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price int64
	Stock int
	gorm.Model
}
