package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/phanthehoang2503/small-project/cart-service/internal/repo"
	"github.com/phanthehoang2503/small-project/internal/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	ProductID uint `json:"product_id" example:"2" binding:"required"`
	Quantity  int  `json:"quantity" example:"3" binding:"required,min=1"`
}

// UpdateQuantityReq struct for update endpoint
type UpdateQuantityReq struct {
	Quantity int `json:"quantity" example:"2" binding:"required"`
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
// @Description Add a product to the cart (can increase quantity if already in the cart).
// @Tags Cart
// @Accept json
// @Produce json
// @Param payload body AddToCartReq true "Add to cart payload"
// @Success 201 {object} handler.CartResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [post]
// @Security BearerAuth
func AddToCart(r *repo.CartRepo, pr *repo.ProductRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in AddToCartReq
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		p, err := pr.Get(in.ProductID)
		if err != nil {
			base := os.Getenv("PRODUCT_SERVICE_URL")
			if base == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "product not found (and service url not set)"})
				return
			}
			url := fmt.Sprintf("%s/%d", base, in.ProductID)
			req, _ := http.NewRequestWithContext(c.Request.Context(), "GET", url, nil)
			otel.GetTextMapPropagator().Inject(c.Request.Context(), propagation.HeaderCarrier(req.Header))
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Failed to fetch product: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Fetch product returned status: %d", resp.StatusCode)
				c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
				return
			}
			defer resp.Body.Close()

			var prod Product
			if err := json.NewDecoder(resp.Body).Decode(&prod); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode product"})
				return
			}

			// Upsert into snapshot
			snapshot := model.ProductSnapshot{
				ProductID: prod.ID,
				Name:      prod.Name,
				Price:     prod.Price,
				Stock:     prod.Stock,
			}
			if err := pr.Upsert(snapshot); err != nil {
				log.Printf("failed to upsert snapshot: %v", err)
			}
			p = &snapshot
		}

		if in.Quantity > p.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock"})
			return
		}

		item := model.Cart{
			UserID:    userID,
			ProductID: p.ProductID,
			Quantity:  in.Quantity,
			Price:     p.Price,
			Subtotal:  p.Price * int64(in.Quantity),
		}

		addedItem, err := r.AddNewItems(&item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, CartResponse{
			ID:        addedItem.ID,
			ProductID: addedItem.ProductID,
			Quantity:  addedItem.Quantity,
			Price:     addedItem.Price,
			Subtotal:  addedItem.Subtotal,
		})
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
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/{id} [put]
// @Security BearerAuth
func UpdateQuantity(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		id := uint(id64)

		var body UpdateQuantityReq
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		updated, err := r.UpdateQuantity(userID, id, body.Quantity)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found in cart"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, CartResponse{
			ID:        updated.ID,
			ProductID: updated.ProductID,
			Quantity:  updated.Quantity,
			Price:     updated.Price,
			Subtotal:  updated.Subtotal,
		})
	}
}

// GetCart godoc
// @Summary List cart items
// @Tags Cart
// @Produce json
// @Success 200 {array} handler.CartResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [get]
// @Security BearerAuth
func GetCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		items, err := r.List(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := make([]CartResponse, 0)
		for _, it := range items {
			resp = append(resp, CartResponse{
				ID:        it.ID,
				ProductID: it.ProductID,
				Quantity:  it.Quantity,
				Price:     it.Price,
				Subtotal:  it.Subtotal,
			})
		}

		c.JSON(http.StatusOK, resp)
	}
}

// RemoveItem godoc
// @Summary Remove an item from cart
// @Tags Cart
// @Param id path int true "Cart item ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cart/{id} [delete]
// @Security BearerAuth
func RemoveItem(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		id := uint(id64)

		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if err := r.Remove(userID, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
				return
			}
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
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart [delete]
// @Security BearerAuth
func ClearCart(r *repo.CartRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if err := r.ClearCart(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
