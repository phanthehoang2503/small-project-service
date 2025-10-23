package model

type Product struct {
	ID    int64  `json: "id"`
	Name  string `json: "name"`
	Price int64  `json: "price"`
}
