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

// struct Product for decode product service response
type Product struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

// AddToCartReq struct used for request body
type AddToCartReq struct {
	ProductID uint `json:"product_id" example:"2"`
	Quantity  int  `json:"quantity" example:"3"`
}

// UpdateQuantityReq struct for update endpoint
type UpdateQuantityReq struct {
	Quantity int `json:"quantity" example: "2"`
}

// CartResponse struct is a public view of cart item
type CartResponse struct {
	ID        uint  `json:"id" example:"1"`
	ProductID uint  `json:"product_id" example:"10"`
	Quantity  int   `json:"quantity" example:"2"`
	Price     int64 `json:"price" example:"10000"`
	Subtotal  int64 `json:"subtotal" example:"20000"`
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to the cart (can add amount of it if already in the cart). Call product-service to get stock and price
// @Tags Cart
// @Accept json
// @Produce json
// @Param payload body AddToCartReq true "Add to cart payload"
// @Success 201 {object} handler.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [post]
func AddToCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in AddToCartReq

		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		base := os.Getenv("PRODUCT_SERVICE_URL")        //e.g: http://localhost:8080/product: url
		url := fmt.Sprintf("%s/%d", base, in.ProductID) //e.g: url/productid
		//Get product info
		resp, err := http.Get(url) //e.g: http://localhost:8080/base/:id (ex: .../product/4)
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

// UpdateQuantity godoc
// @Summary Update cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path int true "Cart item ID"
// @Param payload body UpdateQuantityReq true "New quantity"
// @Success 200 {object} handler.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/{id} [put]
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

// GetCart godoc
// @Summary List cart items
// @Tags Cart
// @Produce json
// @Success 200 {array} handler.CartResponse
// @Failure 500 {object} map[string]string
// @Router /cart [get]
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

// RemoveItem godoc
// @Summary Remove an item from cart
// @Tags Cart
// @Param id path int true "Cart item ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /cart/{id} [delete]
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

// ClearCart godoc
// @Summary Clear all item in the cart
// @Tags Cart
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /cart [delete]
func ClearCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r.ClearCart() != nil {
			c.JSON(http.StatusInternalServerError, r.ClearCart().Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}
