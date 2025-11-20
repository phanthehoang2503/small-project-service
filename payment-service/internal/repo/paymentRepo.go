package repo

import (
	"errors"
	"time"

	"github.com/phanthehoang2503/small-project/payment-service/internal/model"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo { return &PaymentRepo{db: db} }

// Check if payment exists
func (r *PaymentRepo) CreatePending(orderUUID string, amount int64, currency string) (*model.Payment, error) {
	p := &model.Payment{
		OrderUUID: orderUUID,
		Amount:    amount,
		Currency:  currency,
		Provider:  "mock",
		Status:    "PENDING",
	}
	if err := r.db.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PaymentRepo) PaymentSucceeded(orderUUID string) error {
	now := time.Now()
	res := r.db.Model(&model.Payment{}).
		Where("order_uuid = ?", orderUUID).
		Updates(map[string]interface{}{"status": "SUCCEEDED", "updated_at": now})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("payment not found")
	}
	return nil
}

func (r *PaymentRepo) GetByOrderUUID(orderUUID string) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.Where("order_uuid = ?", orderUUID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
