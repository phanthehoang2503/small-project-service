package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"gorm.io/gorm"
)

// UpdateStatusReq binds incoming status update
type UpdateStatusReq struct {
	Status string `json:"status" binding:"required" example:"Paid"`
}

// CreateOrder godoc
// @Summary Create a new order from the current cart
// @Tags Orders
// @Accept json
// @Produce json
// @Param payload body model.Order true "Create order payload"
// @Success 201 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func CreateOrder(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in model.Order
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		base := os.Getenv("CART_SERVICE_URL")

		resp, err := http.Get(base)
		if err != nil || resp.StatusCode != 200 { //<-- better then using to condition check from cart-service
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch cart"})
			return
		}
		defer resp.Body.Close()

		var cartItems []model.OrderItem
		if err := json.NewDecoder(resp.Body).Decode(&cartItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response from cart service"})
			return
		}

		if len(cartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cart empty"})
			return
		}

		order := &model.Order{
			UserId: in.UserId,
			Status: "Pending",
		}

		var total int64
		for _, item := range cartItems {
			order.Items = append(order.Items, model.OrderItem{
				ProductId: item.ID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				Subtotal:  item.Subtotal,
			})
			total += item.Subtotal
		}
		order.Total = total

		if err := r.CreateOrder(order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		created, err := r.GetByID(order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, created)
	}
}

// ListOrders godoc
// @Summary List orders for a user
// @Tags Orders
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {array} model.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func ListOrders(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
			return
		}

		userID64, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		uid := uint(userID64)

		orders, err := r.ListByUser(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

// GetOrder godoc
// @Summary Get order details by id
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.Order
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func GetOrder(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		id := uint(id64)

		order, err := r.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// allowed statuses (simple map)
var allowedStatuses = map[string]bool{
	"Pending":   true,
	"Paid":      true,
	"Cancelled": true,
	"Shipped":   true,
	"Delivered": true,
}

// UpdateOrderStatus godoc
// @Summary Update an order's status
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param payload body UpdateStatusReq true "New status"
// @Success 200 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		orderID := uint(id64)

		var body UpdateStatusReq
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !allowedStatuses[body.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		updated, err := r.UpdateStatus(orderID, body.Status)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}
