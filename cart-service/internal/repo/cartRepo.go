package repo

import (
	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepo struct {
	DB *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{DB: db}
}

func (d *CartRepo) AddUpdateItems(i *model.Cart) (model.Cart, error) { //i: item
	tx := d.DB.Begin()
	if tx.Error != nil {
		return model.Cart{}, tx.Error
	}

	var exist model.Cart
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("product_id = ?", i.ProductID).First(&exist).Error
	if err == nil {
		exist.Quantity += i.Quantity
		exist.Subtotal += exist.Price * int64(exist.Quantity)
		if err := tx.Save(&exist).Error; err != nil {
			tx.Rollback() // -> rollback transaction
			return model.Cart{}, err
		}
	}

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return model.Cart{}, d.DB.Save(&i).Error
	}

	tx.Commit() //-> commit transaction
	return *i, err
} // v.01 change 1: added tx handle to handle added exist item into cart more robust

func (d *CartRepo) List() ([]model.Cart, error) {
	var items []model.Cart
	err := d.DB.Find(&items).Error
	return items, err
}

func (d *CartRepo) UpdateQuantity(id uint, quantity int) (model.Cart, error) {
	var item model.Cart
	if err := d.DB.First(&item, id).Error; err != nil {
		return model.Cart{}, err // gorm.ErrRecordNotFound if missing
	}

	if quantity == 0 {
		if err := d.DB.Delete(&model.Cart{}, id).Error; err != nil {
			return model.Cart{}, err
		}
		return model.Cart{}, gorm.ErrRecordNotFound
	}

	item.Quantity = quantity
	item.Subtotal = item.Price * int64(quantity)

	if err := d.DB.Save(&item).Error; err != nil {
		return model.Cart{}, err
	}
	return item, nil
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
