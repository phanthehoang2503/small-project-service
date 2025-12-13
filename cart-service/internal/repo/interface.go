package repo

import "github.com/phanthehoang2503/small-project/cart-service/internal/model"

type CartRepository interface {
	AddNewItems(i *model.Cart) (model.Cart, error)
	List(UserID uint) ([]model.Cart, error)
	UpdateQuantity(userID, id uint, quantity int) (model.Cart, error)
	Remove(UserID, id uint) error
	ClearCart(userID uint) error
}
