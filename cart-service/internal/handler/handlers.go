package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"gorm.io/gorm"
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

		base := os.Getenv("PRODUCT_SERVICE_URL")        //--> http://localhost:8080/product: url
		url := fmt.Sprintf("%s/%d", base, in.ProductID) //  --> url/productid
		//Get product info
		resp, err := http.Get(url) //--> http://localhost:8080/base/:id (ex: .../product/4)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}

		var p Product
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read product"})
			return
		}
		defer resp.Body.Close()

		if in.Quantity > p.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
			return
		}

		item := model.Cart{
			ProductID: p.ID,
			Quantity:  in.Quantity,
			Price:     p.Price,
			Subtotal:  p.Price * int64(in.Quantity),
		}

		addedItem, err := r.AddUpdateItems(&item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, addedItem)
	}
}

func UpdateQuantity(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var body struct {
			Quantity int `json:"quantity"`
		}
		if c.ShouldBindJSON(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": c.ShouldBindJSON(&body).Error()})
			return
		}

		updated, err := r.UpdateQuantity(uint(id), body.Quantity)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found in cart"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

func GetCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := r.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func RemoveItem(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if err := r.Remove(uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func ClearCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r.ClearCart() != nil {
			c.JSON(http.StatusInternalServerError, r.ClearCart().Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}
