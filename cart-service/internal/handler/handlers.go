package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
)

const productService = "http://localhost:8080/products"

type Product struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

func AddToCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get product info
		resp, err := http.Get(fmt.Sprintf("%s/%d", productService, in.ProductID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var p Product
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read product"})
			return
		}

		if in.Quantity > p.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
			return
		}

		item := model.Cart{
			ProductID: p.ID,
			Quantity:  p.Stock,
			Price:     p.Price,
			Subtotal:  p.Price * int64(in.Quantity),
		}

		if err := r.AddUpdateItems(&item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}
