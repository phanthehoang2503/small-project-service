package repo

import (
	"errors"
	"fmt"

	"github.com/phanthehoang2503/small-project/order-service/internal/model"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Stores order and its items within transaction
func (r *OrderRepo) CreateOrder(userId uint, order *model.Order) error {
	if order == nil || len(order.Items) == 0 {
		return errors.New("invalid order")
	}

	order.UserID = userId

	//compute server-side subtotals and total
	var total int64
	for i := range order.Items {
		item := &order.Items[i]
		if item.Quantity <= 0 {
			return fmt.Errorf("invalid item quantity for product %d", item.ProductID)
		}
		item.Subtotal = int64(item.Quantity) * item.Price
		total += item.Subtotal
	}
	order.Total = total

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Items").Create(order).Error; err != nil {
			return err
		}
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
			if err := tx.Create(&order.Items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OrderRepo) ListByUser(userID uint) ([]model.Order, error) {
	var order []model.Order
	//load related order items link to the user
	/*kinda like: SELECT *
	FROM "orders"
	WHERE "orders"."user_id" = 1 AND "orders"."deleted_at" = NULL
	*/
	if err := r.db.Preload("Items").Where("user_id = ?", userID).Find(&order).Error; err != nil {
		return nil, err
	} /* get all users, and preload all non-cancelled orders
	db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
	*/
	return order, nil
}

func (r *OrderRepo) GetByID(userId, orderId uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Items").
		Where("id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) UpdateStatus(userId, orderId uint, status string) (*model.Order, error) {
	var order model.Order
	if err := r.db.Where("id = ? AND user_id = ?", orderId, userId).First(&order).Error; err != nil {
		return nil, err
	}

	order.Status = status
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Items").First(&order, orderId).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatusByUUID updates order.Status by order.UUID and returns the full order (with items).
// Returns gorm.ErrRecordNotFound if no order matches the UUID.
func (r *OrderRepo) UpdateStatusByUUID(orderUUID, status string) (*model.Order, error) {
	var ord model.Order

	// Update status by uuid
	res := r.db.Model(&model.Order{}).Where("uuid = ?", orderUUID).Update("status", status)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Load and return the full order with items
	if err := r.db.Preload("Items").Where("uuid = ?", orderUUID).First(&ord).Error; err != nil {
		return nil, err
	}
	return &ord, nil
}

// UpdateStatusIfNot updates status only if current status is NOT in forbiddenStates
func (r *OrderRepo) UpdateStatusIfNot(orderUUID, status string, forbiddenStates ...string) (*model.Order, error) {
	if len(forbiddenStates) > 0 {
		var count int64
		// Check if current status is forbidden
		r.db.Model(&model.Order{}).Where("uuid = ? AND status IN ?", orderUUID, forbiddenStates).Count(&count)
		if count > 0 {
			// It is forbidden, return skipping
			return nil, nil // or specific error
		}
	}
	return r.UpdateStatusByUUID(orderUUID, status)
}
