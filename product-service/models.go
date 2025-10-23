package main

import "github.com/google/uuid"

type Product struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Price       int64  `db:"price" json:"price"`
}

func NewProduct(name, desc string, price int64) *Product {
	return &Product{ID: uuid.NewString(), Name: name,
		Description: desc, Price: price}
}
