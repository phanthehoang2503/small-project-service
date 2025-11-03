package model

import "gorm.io/gorm"

type User struct {
	gorm.Model `swaggerignore:"true"`
	Email      string `json:"email" gorm:"uniqueIndex;not null"`
	Username   string `json:"username" gorm:"uniqueIndex;not null"`
	Password   string `json:"-"` // hashed password
}
