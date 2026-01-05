package model

import "gorm.io/gorm"

type User struct {
	gorm.Model `swaggerignore:"true"`
	Email      string `json:"email" gorm:"uniqueIndex;not null" example:"user@example.com"`
	Username   string `json:"username" gorm:"uniqueIndex;not null" example:"username123"`
	Password   string `json:"-"` // hashed password
}
