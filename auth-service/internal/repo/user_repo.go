package repo

import (
	"errors"

	"github.com/phanthehoang2503/small-project/auth-service/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("user not found")

type UserRepo interface {
	Create(u *model.User) error
	GetUser(value string) (*model.User, error)
}

type userRepoDB struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepoDB{db: db}
}

func (r *userRepoDB) Create(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *userRepoDB) GetUser(value string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ? OR username = ?", value, value).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
