package repository

import (
	"context"
	"errors"
	"fmt"

	"bicycle-rent-api/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type voucherRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewVoucherRepositoryPostgres(db *pgxpool.Pool) VoucherRepository {
	return &voucherRepositoryPostgres{db: db}
}

// =========================
// Create
// =========================
func (r *voucherRepositoryPostgres) Create(v entity.Voucher) (entity.Voucher, error) {
	query := `
        INSERT INTO vouchers (
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING
            voucher_id,
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active,
            created_at
    `

	err := r.db.QueryRow(
		context.Background(),
		query,
		v.VoucherCode,
		v.DiscountPercentage,
		v.MaxDiscountAmount,
		v.ValidUntil,
		v.IsActive,
	).Scan(
		&v.VoucherID,
		&v.VoucherCode,
		&v.DiscountPercentage,
		&v.MaxDiscountAmount,
		&v.ValidUntil,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err != nil {
		return entity.Voucher{}, fmt.Errorf("failed to create voucher: %w", err)
	}
	return v, nil
}

// =========================
// Get All
// =========================
func (r *voucherRepositoryPostgres) GetAll() ([]entity.Voucher, error) {
	query := `
        SELECT
            voucher_id,
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active,
            created_at
        FROM vouchers
        ORDER BY voucher_id
    `

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get vouchers: %w", err)
	}
	defer rows.Close()

	var vouchers []entity.Voucher
	for rows.Next() {
		var v entity.Voucher
		if err := rows.Scan(
			&v.VoucherID,
			&v.VoucherCode,
			&v.DiscountPercentage,
			&v.MaxDiscountAmount,
			&v.ValidUntil,
			&v.IsActive,
			&v.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan voucher: %w", err)
		}
		vouchers = append(vouchers, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read vouchers: %w", err)
	}
	return vouchers, nil
}

// =========================
// Get By ID
// =========================
func (r *voucherRepositoryPostgres) GetByID(id int64) (entity.Voucher, error) {
	query := `
        SELECT
            voucher_id,
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active,
            created_at
        FROM vouchers
        WHERE voucher_id = $1
    `

	var v entity.Voucher
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&v.VoucherID,
		&v.VoucherCode,
		&v.DiscountPercentage,
		&v.MaxDiscountAmount,
		&v.ValidUntil,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Voucher{}, errors.New("voucher not found")
		}
		return entity.Voucher{}, fmt.Errorf("failed to get voucher: %w", err)
	}
	return v, nil
}

// =========================
// Get By Code
// =========================
func (r *voucherRepositoryPostgres) GetByCode(code string) (entity.Voucher, error) {
	query := `
        SELECT
            voucher_id,
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active,
            created_at
        FROM vouchers
        WHERE voucher_code = $1
    `

	var v entity.Voucher
	err := r.db.QueryRow(context.Background(), query, code).Scan(
		&v.VoucherID,
		&v.VoucherCode,
		&v.DiscountPercentage,
		&v.MaxDiscountAmount,
		&v.ValidUntil,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Voucher{}, errors.New("voucher not found")
		}
		return entity.Voucher{}, fmt.Errorf("failed to get voucher by code: %w", err)
	}
	return v, nil
}

// =========================
// Update
// =========================
func (r *voucherRepositoryPostgres) Update(v entity.Voucher) (entity.Voucher, error) {
	query := `
        UPDATE vouchers
        SET
            voucher_code = $1,
            discount_percentage = $2,
            max_discount_amount = $3,
            valid_until = $4,
            is_active = $5
        WHERE voucher_id = $6
        RETURNING
            voucher_id,
            voucher_code,
            discount_percentage,
            max_discount_amount,
            valid_until,
            is_active,
            created_at
    `

	err := r.db.QueryRow(
		context.Background(),
		query,
		v.VoucherCode,
		v.DiscountPercentage,
		v.MaxDiscountAmount,
		v.ValidUntil,
		v.IsActive,
		v.VoucherID,
	).Scan(
		&v.VoucherID,
		&v.VoucherCode,
		&v.DiscountPercentage,
		&v.MaxDiscountAmount,
		&v.ValidUntil,
		&v.IsActive,
		&v.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Voucher{}, errors.New("voucher not found")
		}
		return entity.Voucher{}, fmt.Errorf("failed to update voucher: %w", err)
	}
	return v, nil
}

// =========================
// Delete
// =========================
func (r *voucherRepositoryPostgres) Delete(id int64) error {
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM vouchers WHERE voucher_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete voucher: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("voucher not found")
	}
	return nil
}
