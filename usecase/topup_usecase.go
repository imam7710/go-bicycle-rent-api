package usecase

import (
	"errors"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/repository"
)

type TopupUsecase struct {
	topupRepository repository.TopupRepository
}

func NewTopupUsecase(topupRepository repository.TopupRepository) *TopupUsecase {
	return &TopupUsecase{
		topupRepository: topupRepository,
	}
}

func (u *TopupUsecase) Create(userID int64, amount float64) (entity.Topup, error) {
	if userID <= 0 {
		return entity.Topup{}, errors.New("invalid user")
	}

	if amount <= 0 {
		return entity.Topup{}, errors.New("top up amount must be greater than 0")
	}

	topup := entity.Topup{
		UserID: userID,
		Amount: amount,
	}

	createdTopup, err := u.topupRepository.Create(topup)
	if err != nil {
		return entity.Topup{}, err
	}

	return createdTopup, nil
}

func (u *TopupUsecase) GetByUserID(userID int64) ([]entity.Topup, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user")
	}

	return u.topupRepository.GetByUserID(userID)
}
