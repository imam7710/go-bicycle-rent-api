package usecase

import (
	"errors"
	"strings"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/repository"
)

type BicycleUsecase struct {
	bicycleRepository repository.BicycleRepository
}

func NewBicycleUsecase(bicycleRepository repository.BicycleRepository) *BicycleUsecase {
	return &BicycleUsecase{
		bicycleRepository: bicycleRepository,
	}
}

func (u *BicycleUsecase) Create(bicycle entity.Bicycle) (entity.Bicycle, error) {
	bicycle.BicycleModel = strings.TrimSpace(bicycle.BicycleModel)
	bicycle.BicycleType = strings.TrimSpace(bicycle.BicycleType)

	if bicycle.BicycleModel == "" {
		return entity.Bicycle{}, errors.New("bicycle model is required")
	}

	if bicycle.BicycleType == "" {
		return entity.Bicycle{}, errors.New("bicycle type is required")
	}

	if bicycle.BicycleStock < 0 {
		return entity.Bicycle{}, errors.New("bicycle stock cannot be negative")
	}

	if bicycle.PricePerDayNormal <= 0 {
		return entity.Bicycle{}, errors.New("normal daily price must be greater than 0")
	}

	if bicycle.PricePerDayHoliday <= 0 {
		return entity.Bicycle{}, errors.New("holiday daily price must be greater than 0")
	}

	return u.bicycleRepository.Create(bicycle)
}

func (u *BicycleUsecase) GetAll() ([]entity.Bicycle, error) {
	return u.bicycleRepository.GetAll()
}

func (u *BicycleUsecase) GetByID(bicycleID int64) (entity.Bicycle, error) {
	if bicycleID <= 0 {
		return entity.Bicycle{}, errors.New("invalid bicycle id")
	}

	return u.bicycleRepository.GetByID(bicycleID)
}

func (u *BicycleUsecase) Update(bicycle entity.Bicycle) (entity.Bicycle, error) {
	if bicycle.BicycleID <= 0 {
		return entity.Bicycle{}, errors.New("invalid bicycle id")
	}

	if bicycle.BicycleModel == "" {
		return entity.Bicycle{}, errors.New("bicycle model is required")
	}

	if bicycle.BicycleType == "" {
		return entity.Bicycle{}, errors.New("bicycle type is required")
	}

	if bicycle.BicycleStock < 0 {
		return entity.Bicycle{}, errors.New("bicycle stock cannot be negative")
	}

	if bicycle.PricePerDayNormal <= 0 {
		return entity.Bicycle{}, errors.New("normal daily price must be greater than 0")
	}

	if bicycle.PricePerDayHoliday <= 0 {
		return entity.Bicycle{}, errors.New("holiday daily price must be greater than 0")
	}

	return u.bicycleRepository.Update(bicycle)
}

func (u *BicycleUsecase) Delete(bicycleID int64) error {
	if bicycleID <= 0 {
		return errors.New("invalid bicycle id")
	}

	return u.bicycleRepository.Delete(bicycleID)
}
