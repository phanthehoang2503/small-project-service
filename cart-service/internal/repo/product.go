package repo

import (
	"time"

	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Upsert(s model.ProductSnapshot) error {
	s.UpdatedAt = time.Now().UTC()
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "price", "stock", "version", "updated_at"}),
	}).Create(&s).Error
}

func (r *ProductRepo) Delete(id uint) error {
	return r.db.Delete(&model.ProductSnapshot{}, id).Error
}

func (r *ProductRepo) Get(id uint) (*model.ProductSnapshot, error) {
	var p model.ProductSnapshot
	err := r.db.First(&p, id).Error
	return &p, err
}
