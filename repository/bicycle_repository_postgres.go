package repository

import (
	"context"
	"errors"
	"fmt"

	"bicycle-rent-api/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bicycleRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewBicycleRepositoryPostgres(db *pgxpool.Pool) BicycleRepository {
	return &bicycleRepositoryPostgres{
		db: db,
	}
}

func (r *bicycleRepositoryPostgres) Create(bicycle entity.Bicycle) (entity.Bicycle, error) {
	query := `
		INSERT INTO bicycles (
			bicycle_model,
			bicycle_type,
			bicycle_stock,
			price_per_day_normal,
			price_per_day_holiday
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			bicycle_id,
			bicycle_model,
			bicycle_type,
			bicycle_stock,
			price_per_day_normal,
			price_per_day_holiday,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		bicycle.BicycleModel,
		bicycle.BicycleType,
		bicycle.BicycleStock,
		bicycle.PricePerDayNormal,
		bicycle.PricePerDayHoliday,
	).Scan(
		&bicycle.BicycleID,
		&bicycle.BicycleModel,
		&bicycle.BicycleType,
		&bicycle.BicycleStock,
		&bicycle.PricePerDayNormal,
		&bicycle.PricePerDayHoliday,
		&bicycle.CreatedAt,
		&bicycle.UpdatedAt,
	)

	if err != nil {
		return entity.Bicycle{}, fmt.Errorf("failed to create bicycle: %w", err)
	}

	return bicycle, nil
}

func (r *bicycleRepositoryPostgres) GetAll() ([]entity.Bicycle, error) {
	query := `
		SELECT
			bicycle_id,
			bicycle_model,
			bicycle_type,
			bicycle_stock,
			price_per_day_normal,
			price_per_day_holiday,
			created_at,
			updated_at
		FROM bicycles
		ORDER BY bicycle_id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get bicycles: %w", err)
	}
	defer rows.Close()

	var bicycles []entity.Bicycle

	for rows.Next() {
		var bicycle entity.Bicycle

		err := rows.Scan(
			&bicycle.BicycleID,
			&bicycle.BicycleModel,
			&bicycle.BicycleType,
			&bicycle.BicycleStock,
			&bicycle.PricePerDayNormal,
			&bicycle.PricePerDayHoliday,
			&bicycle.CreatedAt,
			&bicycle.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan bicycle: %w", err)
		}

		bicycles = append(bicycles, bicycle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read bicycles: %w", err)
	}

	return bicycles, nil
}

func (r *bicycleRepositoryPostgres) GetByID(bicycleID int64) (entity.Bicycle, error) {
	query := `
		SELECT
			bicycle_id,
			bicycle_model,
			bicycle_type,
			bicycle_stock,
			price_per_day_normal,
			price_per_day_holiday,
			created_at,
			updated_at
		FROM bicycles
		WHERE bicycle_id = $1
	`

	var bicycle entity.Bicycle

	err := r.db.QueryRow(
		context.Background(),
		query,
		bicycleID,
	).Scan(
		&bicycle.BicycleID,
		&bicycle.BicycleModel,
		&bicycle.BicycleType,
		&bicycle.BicycleStock,
		&bicycle.PricePerDayNormal,
		&bicycle.PricePerDayHoliday,
		&bicycle.CreatedAt,
		&bicycle.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bicycle{}, errors.New("bicycle not found")
		}

		return entity.Bicycle{}, fmt.Errorf("failed to get bicycle: %w", err)
	}

	return bicycle, nil
}

func (r *bicycleRepositoryPostgres) Update(bicycle entity.Bicycle) (entity.Bicycle, error) {
	query := `
		UPDATE bicycles
		SET
			bicycle_model = $1,
			bicycle_type = $2,
			bicycle_stock = $3,
			price_per_day_normal = $4,
			price_per_day_holiday = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE bicycle_id = $6
		RETURNING
			bicycle_id,
			bicycle_model,
			bicycle_type,
			bicycle_stock,
			price_per_day_normal,
			price_per_day_holiday,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		bicycle.BicycleModel,
		bicycle.BicycleType,
		bicycle.BicycleStock,
		bicycle.PricePerDayNormal,
		bicycle.PricePerDayHoliday,
		bicycle.BicycleID,
	).Scan(
		&bicycle.BicycleID,
		&bicycle.BicycleModel,
		&bicycle.BicycleType,
		&bicycle.BicycleStock,
		&bicycle.PricePerDayNormal,
		&bicycle.PricePerDayHoliday,
		&bicycle.CreatedAt,
		&bicycle.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Bicycle{}, errors.New("bicycle not found")
		}

		return entity.Bicycle{}, fmt.Errorf("failed to update bicycle: %w", err)
	}

	return bicycle, nil
}

func (r *bicycleRepositoryPostgres) Delete(bicycleID int64) error {
	result, err := r.db.Exec(
		context.Background(),
		`DELETE FROM bicycles WHERE bicycle_id = $1`,
		bicycleID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete bicycle: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("bicycle not found")
	}

	return nil
}

func (r *bicycleRepositoryPostgres) UpdateStock(bicycleID int64, stock int) error {
	result, err := r.db.Exec(
		context.Background(),
		`
		UPDATE bicycles
		SET
			bicycle_stock = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE bicycle_id = $2
		`,
		stock,
		bicycleID,
	)

	if err != nil {
		return fmt.Errorf("failed to update bicycle stock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("bicycle not found")
	}

	return nil
}
