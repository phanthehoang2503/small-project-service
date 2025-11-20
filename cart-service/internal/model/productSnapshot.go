package model

import "time"

type ProductSnapshot struct {
	ProductID uint      `gorm:"primaryKey;column:product_id" json:"product_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	UpdatedAt time.Time `json:"updated_at"`
}
