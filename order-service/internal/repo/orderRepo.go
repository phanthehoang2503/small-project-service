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
func (r *OrderRepo) CreateOrder(orders *model.Order) error {
	if orders == nil || len(orders.Item) == 0 {
		return errors.New("invalid order")
	}
	// same as func (...) AddUpdateItems from cartRepo in cart-service but in a cleaner way.
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(orders).Error; err != nil {
			return err
		} // check if ini order yet
		for i := range orders.Item {
			orders.Item[i].OrderId = orders.ID //assign key
			if err := tx.Create(&orders.Item[i]).Error; err != nil {
				return err
			}
		}
		return nil
	}) // return error if any line in this got error and revoke all the transaction or vice versa
}
