package repository

import "bicycle-rent-api/entity"

type TopupRepository interface {
	Create(topup entity.Topup) (entity.Topup, error)
	GetByUserID(userID int64) ([]entity.Topup, error)
}
