package repository

import (
	"context"
	"errors"
	"fmt"

	"bicycle-rent-api/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) UserRepository {
	return &userRepositoryPostgres{
		db: db,
	}
}

func (r *userRepositoryPostgres) Create(user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (
			user_name,
			email,
			password_hash,
			role,
			balance
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			user_id,
			user_name,
			email,
			password_hash,
			role,
			balance,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.UserName,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Balance,
	).Scan(
		&user.UserID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *userRepositoryPostgres) GetByID(userID int64) (entity.User, error) {
	query := `
		SELECT
			user_id,
			user_name,
			email,
			password_hash,
			role,
			balance,
			created_at,
			updated_at
		FROM users
		WHERE user_id = $1
	`

	var user entity.User

	err := r.db.QueryRow(
		context.Background(),
		query,
		userID,
	).Scan(
		&user.UserID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}

		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepositoryPostgres) GetByEmail(email string) (entity.User, error) {
	query := `
		SELECT
			user_id,
			user_name,
			email,
			password_hash,
			role,
			balance,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var user entity.User

	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&user.UserID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}

		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepositoryPostgres) Update(user entity.User) (entity.User, error) {
	query := `
		UPDATE users
		SET
			user_name = $1,
			email = $2,
			role = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $4
		RETURNING
			user_id,
			user_name,
			email,
			password_hash,
			role,
			balance,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.UserName,
		user.Email,
		user.Role,
		user.UserID,
	).Scan(
		&user.UserID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}

		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *userRepositoryPostgres) UpdateBalance(userID int64, amount float64) error {
	query := `
		UPDATE users
		SET
			balance = balance + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
		AND balance + $1 >= 0
	`

	result, err := r.db.Exec(
		context.Background(),
		query,
		amount,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or insufficient balance")
	}

	return nil
}
