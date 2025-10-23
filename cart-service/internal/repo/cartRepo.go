package repo

import (
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"gorm.io/gorm"
)

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{DB: db}
}

func (d *CartRepo) AddUpdateItems(i *model.Cart) error { //i: item
	var exist model.Cart
	err := d.DB.Where("product_id = ?", i.ProductID).First(&exist).Error
	if err == nil {
		exist.Quantity += i.Quantity
		exist.Subtotal += exist.Price * int64(exist.Quantity)
		return d.DB.Save(&exist).Error
	}

	if err == gorm.ErrRecordNotFound {
		return d.DB.Save(&i).Error
	}

	return err
}

func (d *CartRepo) List() ([]model.Cart, error) {
	var items []model.Cart
	err := d.DB.Find(&items).Error
	return items, err
}

func (d *CartRepo) UpdateQuantity(id uint, quantity int) error {
	var item model.Cart
	if err := d.DB.First(&item, id).Error; err != nil {
		return err
	}
	item.Quantity = quantity
	item.Subtotal = item.Price * int64(quantity)
	return d.DB.Save(&item).Error
}

func (d *CartRepo) Remove(id uint) error {
	res := d.DB.Delete(&model.Cart{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *CartRepo) ClearCart() error {
	return d.DB.Exec("DELETE FROM cart_items").Error
}
