package repository

import "bicycle-rent-api/entity"

type BicycleRepository interface {
	Create(bicycle entity.Bicycle) (entity.Bicycle, error)
	GetAll() ([]entity.Bicycle, error)
	GetByID(bicycleID int64) (entity.Bicycle, error)
	Update(bicycle entity.Bicycle) (entity.Bicycle, error)
	Delete(bicycleID int64) error
	UpdateStock(bicycleID int64, stock int) error
}
