package repo

import (
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	DB *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (r *ProductRepo) Upsert(p model.ProductSnapshot) error {
	return r.DB.
		Where("id = ?", p.ID).
		Assign(p).
		FirstOrCreate(&p).Error
}

func (r *ProductRepo) Delete(id uint) error {
	return r.DB.Delete(&model.ProductSnapshot{}, id).Error
}

func (r *ProductRepo) Get(id uint) (*model.ProductSnapshot, error) {
	var p model.ProductSnapshot
	err := r.DB.First(&p, id).Error
	return &p, err
}
