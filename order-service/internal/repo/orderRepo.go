package repo

import (
	"errors"

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
func (r *OrderRepo) CreateOrder(order *model.Order) error {
	if order == nil || len(order.Items) == 0 {
		return errors.New("invalid order")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
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

func (r *OrderRepo) GetByID(orderId uint) (*model.Order, error) {
	var order model.Order
	if err := r.db.Preload("Items").First(&order, orderId).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) UpdateStatus(orderId uint, status string) (*model.Order, error) {
	var order model.Order
	if err := r.db.First(&order, orderId).Error; err != nil {
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
