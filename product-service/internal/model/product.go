package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name" example:"Smartphone"`
	Price      int64  `json:"price" example:"999"`
	Stock      int    `json:"stock" example:"100"`
}
