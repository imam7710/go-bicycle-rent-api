package repository

import (
	"context"
	"errors"
	"fmt"

	"bicycle-rent-api/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type topupRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewTopupRepositoryPostgres(db *pgxpool.Pool) TopupRepository {
	return &topupRepositoryPostgres{
		db: db,
	}
}

func (r *topupRepositoryPostgres) Create(topup entity.Topup) (entity.Topup, error) {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Topup{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	// Update user balance
	var balance float64

	err = tx.QueryRow(
		ctx,
		`
		UPDATE users
		SET
			balance = balance + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
		RETURNING balance
		`,
		topup.Amount,
		topup.UserID,
	).Scan(&balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Topup{}, fmt.Errorf("user not found")
		}

		return entity.Topup{}, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create topup history
	query := `
		INSERT INTO topups (
			user_id,
			amount
		)
		VALUES ($1, $2)
		RETURNING
			topup_id,
			user_id,
			amount,
			created_at
	`

	err = tx.QueryRow(
		ctx,
		query,
		topup.UserID,
		topup.Amount,
	).Scan(
		&topup.TopupID,
		&topup.UserID,
		&topup.Amount,
		&topup.CreatedAt,
	)

	if err != nil {
		return entity.Topup{}, fmt.Errorf("failed to create topup: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Topup{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return topup, nil
}

func (r *topupRepositoryPostgres) GetByUserID(userID int64) ([]entity.Topup, error) {
	query := `
		SELECT
			topup_id,
			user_id,
			amount,
			created_at
		FROM topups
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get topups: %w", err)
	}
	defer rows.Close()

	var topups []entity.Topup

	for rows.Next() {
		var topup entity.Topup

		err := rows.Scan(
			&topup.TopupID,
			&topup.UserID,
			&topup.Amount,
			&topup.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan topup: %w", err)
		}

		topups = append(topups, topup)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read topups: %w", err)
	}

	return topups, nil
}
