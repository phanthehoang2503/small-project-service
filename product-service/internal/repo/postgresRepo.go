package repo

import (
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func (d *Database) Create(p model.Product) (model.Product, error) {
	if err := d.DB.Create(&p).Error; err != nil {
		return model.Product{}, err
	}
	return p, nil
}

func (d *Database) List() ([]model.Product, error) {
	var products []model.Product
	if err := d.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (d *Database) Get(id int64) (model.Product, error) {
	var p model.Product

	if err := d.DB.First(&p, id).Error; err != nil {
		return model.Product{}, err
	}
	return p, nil
}

func (d *Database) Update(id int64, newData model.Product) (model.Product, error) {
	var exist model.Product
	if err := d.DB.First(&exist, id).Error; err != nil {
		return model.Product{}, err
	}

	exist.Name = newData.Name
	exist.Price = newData.Price
	exist.Stock = newData.Stock
	if err := d.DB.Save(&exist).Error; err != nil {
		return model.Product{}, err
	}
	return exist, nil
}

func (d *Database) Delete(id int64) error {
	res := d.DB.Delete(&model.Product{}, id)
	if res.Error != nil { //err while delete
		return res.Error
	}
	if res.RowsAffected == 0 { //not found record
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *Database) DeductStock(productID uint, quantity int) error {
	res := d.DB.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
