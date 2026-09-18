package usecase

import (
	"errors"
	"strings"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) Register(user entity.User) (entity.User, error) {
	user.UserName = strings.TrimSpace(user.UserName)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if user.UserName == "" {
		return entity.User{}, errors.New("name is required")
	}

	if user.Email == "" {
		return entity.User{}, errors.New("email is required")
	}

	if user.PasswordHash == "" {
		return entity.User{}, errors.New("password is required")
	}

	if len(user.PasswordHash) < 6 {
		return entity.User{}, errors.New("password must be at least 6 characters")
	}

	_, err := u.userRepository.GetByEmail(user.Email)
	if err == nil {
		return entity.User{}, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.PasswordHash),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return entity.User{}, errors.New("failed to hash password")
	}

	user.PasswordHash = string(hashedPassword)
	user.Role = "customer"
	user.Balance = 0

	createdUser, err := u.userRepository.Create(user)
	if err != nil {
		return entity.User{}, err
	}

	return createdUser, nil
}

func (u *UserUsecase) Login(email, password string) (entity.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return entity.User{}, errors.New("email and password are required")
	}

	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		return entity.User{}, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return entity.User{}, errors.New("invalid email or password")
	}

	return user, nil
}

func (u *UserUsecase) GetByID(userID int64) (entity.User, error) {
	user, err := u.userRepository.GetByID(userID)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
