package repository

import "bicycle-rent-api/entity"

type UserRepository interface {
	Create(user entity.User) (entity.User, error)
	GetByID(userID int64) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	UpdateBalance(userID int64, amount float64) error
}
