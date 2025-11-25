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
	"github.com/google/uuid"
	"github.com/phanthehoang2503/small-project/internal/broker"
	"github.com/phanthehoang2503/small-project/internal/logger"
	"github.com/phanthehoang2503/small-project/internal/util"
	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"github.com/phanthehoang2503/small-project/order-service/internal/publisher"
	"github.com/phanthehoang2503/small-project/order-service/internal/repo"
	"gorm.io/gorm"
)

// UpdateStatusReq binds incoming status update
type UpdateStatusReq struct {
	Status string `json:"status" binding:"required" example:"Paid"`
}

type cartItemResp struct {
	CartID    uint  `json:"ID,omitempty"`
	UserID    uint  `json:"user_id"`
	ProductID uint  `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
	Subtotal  int64 `json:"subtotal"`
}

// CreateOrder godoc
// @Summary Create a new order from the current cart
// @Tags Orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func CreateOrder(r *repo.OrderRepo, b *broker.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user from JWT in context
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		var in model.Order
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		base := os.Getenv("CART_SERVICE_URL")
		if base == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cart service url not configured"})
			return
		}

		cartURL := base
		req, err := http.NewRequest("GET", cartURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}
		req.Header.Set("Authorization", c.GetHeader("Authorization"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch cart"})
			return
		}
		defer resp.Body.Close()

		var cartItems []struct {
			ProductID uint  `json:"product_id"`
			Quantity  int   `json:"quantity"`
			Price     int64 `json:"price"`
			Subtotal  int64 `json:"subtotal"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&cartItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid response from cart service"})
			return
		}

		if len(cartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cart empty"})
			return
		}

		order := &model.Order{
			UserID: userID,
			UUID:   uuid.New().String(),
			Status: "Pending",
		}

		var total int64
		for _, item := range cartItems {
			if item.Quantity <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item quantity"})
				return
			}
			oi := model.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
			oi.Subtotal = int64(oi.Quantity) * oi.Price
			order.Items = append(order.Items, oi)
			total += oi.Subtotal
		}
		order.Total = total

		if err := r.CreateOrder(userID, order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		created, err := r.GetByID(userID, order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if b != nil {
			if err := publisher.PublishOrderRequested(b, created.UUID, created.UUID, created.UserID, created.Total, "VND"); err != nil {
				log.Printf("failed to publish order requested event: %v", err)
			}
		}

		c.JSON(http.StatusCreated, created)

		logger.Info(c.Request.Context(), fmt.Sprintf("Order created: id=%d uuid=%s user_id=%d total=%d", created.ID, created.UUID, created.UserID, created.Total))
	}
}

// ListOrders godoc
// @Summary List orders for the authenticated user
// @Tags Orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func ListOrders(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		orders, err := r.ListByUser(userID)
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
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.Order
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func GetOrder(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		idStr := c.Param("id")
		id64, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		id := uint(id64)

		order, err := r.GetByID(userID, id)
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

// allowed statuses ( simple map )
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param payload body UpdateStatusReq true "New status"
// @Success 200 {object} model.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(r *repo.OrderRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

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

		updated, err := r.UpdateStatus(userID, orderID, body.Status)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)

		logger.Info(c.Request.Context(), fmt.Sprintf("Order status updated: id=%d status=%s", updated.ID, updated.Status))
	}
}
