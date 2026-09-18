package repository

import (
	"context"
	"errors"
	"fmt"

	"bicycle-rent-api/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookingRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewBookingRepositoryPostgres(db *pgxpool.Pool) BookingRepository {
	return &bookingRepositoryPostgres{
		db: db,
	}
}

func (r *bookingRepositoryPostgres) Create(booking entity.Booking) (entity.Booking, error) {

	query := `
		INSERT INTO bookings (
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
		RETURNING
			booking_id,
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		booking.UserID,
		booking.BicycleID,
		booking.VoucherID,
		booking.BookingDate,
		booking.BookingStatus,
		booking.RentalType,
		booking.Duration,
		booking.TimeDepart,
		booking.TimeArrive,
		booking.SubtotalAmount,
		booking.DiscountAmount,
		booking.TotalAmount,
	).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.BicycleID,
		&booking.VoucherID,
		&booking.BookingDate,
		&booking.BookingStatus,
		&booking.RentalType,
		&booking.Duration,
		&booking.TimeDepart,
		&booking.TimeArrive,
		&booking.SubtotalAmount,
		&booking.DiscountAmount,
		&booking.TotalAmount,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to create booking: %w",
			err,
		)
	}

	return booking, nil
}

func (r *bookingRepositoryPostgres) GetByID(bookingID int64) (entity.Booking, error) {

	query := `
		SELECT
			booking_id,
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount,
			created_at,
			updated_at
		FROM bookings
		WHERE booking_id = $1
	`

	var booking entity.Booking

	err := r.db.QueryRow(
		context.Background(),
		query,
		bookingID,
	).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.BicycleID,
		&booking.VoucherID,
		&booking.BookingDate,
		&booking.BookingStatus,
		&booking.RentalType,
		&booking.Duration,
		&booking.TimeDepart,
		&booking.TimeArrive,
		&booking.SubtotalAmount,
		&booking.DiscountAmount,
		&booking.TotalAmount,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Booking{}, errors.New("booking not found")
		}

		return entity.Booking{}, fmt.Errorf(
			"failed to get booking: %w",
			err,
		)
	}

	return booking, nil
}

func (r *bookingRepositoryPostgres) GetByUserID(userID int64) ([]entity.Booking, error) {

	query := `
		SELECT
			booking_id,
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount,
			created_at,
			updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY booking_date DESC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get user bookings: %w",
			err,
		)
	}

	defer rows.Close()

	var bookings []entity.Booking

	for rows.Next() {
		var booking entity.Booking

		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.BicycleID,
			&booking.VoucherID,
			&booking.BookingDate,
			&booking.BookingStatus,
			&booking.RentalType,
			&booking.Duration,
			&booking.TimeDepart,
			&booking.TimeArrive,
			&booking.SubtotalAmount,
			&booking.DiscountAmount,
			&booking.TotalAmount,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan booking: %w",
				err,
			)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to read bookings: %w",
			err,
		)
	}

	return bookings, nil
}

func (r *bookingRepositoryPostgres) UpdateStatus(bookingID int64, status string) error {

	query := `
		UPDATE bookings
		SET
			booking_status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE booking_id = $2
	`

	result, err := r.db.Exec(
		context.Background(),
		query,
		status,
		bookingID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to update booking status: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return errors.New("booking not found")
	}

	return nil
}

func (r *bookingRepositoryPostgres) GetReport() (BookingReport, error) {

	query := `
		SELECT
			COUNT(*) AS total_bookings,
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) FILTER (
				WHERE booking_status = 'completed'
			) AS completed_bookings,
			COUNT(*) FILTER (
				WHERE booking_status = 'cancelled'
			) AS cancelled_bookings
		FROM bookings
	`

	var report BookingReport

	err := r.db.QueryRow(
		context.Background(),
		query,
	).Scan(
		&report.TotalBookings,
		&report.TotalRevenue,
		&report.CompletedBookings,
		&report.CancelledBookings,
	)

	if err != nil {
		return BookingReport{}, fmt.Errorf(
			"failed to get booking report: %w",
			err,
		)
	}

	return report, nil
}

func (r *bookingRepositoryPostgres) CreateTransaction(booking entity.Booking, totalAmount float64) (entity.Booking, error) {

	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	// 1. Kurangi saldo user
	var currentBalance float64

	err = tx.QueryRow(
		ctx,
		`
		SELECT balance
		FROM users
		WHERE user_id = $1
		FOR UPDATE
		`,
		booking.UserID,
	).Scan(&currentBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Booking{}, errors.New("user not found")
		}

		return entity.Booking{}, fmt.Errorf(
			"failed to get user balance: %w",
			err,
		)
	}

	if currentBalance < totalAmount {
		return entity.Booking{}, errors.New("insufficient balance")
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE users
		SET
			balance = balance - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
		`,
		totalAmount,
		booking.UserID,
	)

	if err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to update user balance: %w",
			err,
		)
	}

	// 2. Kurangi stock bicycle
	var currentStock int

	err = tx.QueryRow(
		ctx,
		`
		SELECT bicycle_stock
		FROM bicycles
		WHERE bicycle_id = $1
		FOR UPDATE
		`,
		booking.BicycleID,
	).Scan(&currentStock)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Booking{}, errors.New("bicycle not found")
		}

		return entity.Booking{}, fmt.Errorf(
			"failed to get bicycle stock: %w",
			err,
		)
	}

	if currentStock <= 0 {
		return entity.Booking{}, errors.New("bicycle is out of stock")
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE bicycles
		SET
			bicycle_stock = bicycle_stock - 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE bicycle_id = $1
		`,
		booking.BicycleID,
	)

	if err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to update bicycle stock: %w",
			err,
		)
	}

	// 3. Simpan booking
	booking.TotalAmount = totalAmount

	query := `
		INSERT INTO bookings (
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12
		)
		RETURNING
			booking_id,
			user_id,
			bicycle_id,
			voucher_id,
			booking_date,
			booking_status,
			rental_type,
			duration,
			time_depart,
			time_arrive,
			subtotal_amount,
			discount_amount,
			total_amount,
			created_at,
			updated_at
	`

	err = tx.QueryRow(
		ctx,
		query,
		booking.UserID,
		booking.BicycleID,
		booking.VoucherID,
		booking.BookingDate,
		booking.BookingStatus,
		booking.RentalType,
		booking.Duration,
		booking.TimeDepart,
		booking.TimeArrive,
		booking.SubtotalAmount,
		booking.DiscountAmount,
		booking.TotalAmount,
	).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.BicycleID,
		&booking.VoucherID,
		&booking.BookingDate,
		&booking.BookingStatus,
		&booking.RentalType,
		&booking.Duration,
		&booking.TimeDepart,
		&booking.TimeArrive,
		&booking.SubtotalAmount,
		&booking.DiscountAmount,
		&booking.TotalAmount,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to create booking: %w",
			err,
		)
	}

	// 4. Commit
	if err := tx.Commit(ctx); err != nil {
		return entity.Booking{}, fmt.Errorf(
			"failed to commit transaction: %w",
			err,
		)
	}

	return booking, nil
}

func (r *bookingRepositoryPostgres) Cancel(bookingID int64) error {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to begin cancellation transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	// Get booking
	var (
		userID    int64
		bicycleID int64
		total     float64
		status    string
	)

	err = tx.QueryRow(
		ctx,
		`
		SELECT
			user_id,
			bicycle_id,
			total_amount,
			booking_status
		FROM bookings
		WHERE booking_id = $1
		FOR UPDATE
		`,
		bookingID,
	).Scan(
		&userID,
		&bicycleID,
		&total,
		&status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("booking not found")
		}

		return fmt.Errorf(
			"failed to get booking: %w",
			err,
		)
	}

	// Only confirmed booking can be cancelled
	if status != "confirmed" {
		return errors.New(
			"only confirmed booking can be cancelled",
		)
	}

	// Return balance
	_, err = tx.Exec(
		ctx,
		`
		UPDATE users
		SET
			balance = balance + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
		`,
		total,
		userID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to refund user balance: %w",
			err,
		)
	}

	// Return bicycle stock
	_, err = tx.Exec(
		ctx,
		`
		UPDATE bicycles
		SET
			bicycle_stock = bicycle_stock + 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE bicycle_id = $1
		`,
		bicycleID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to restore bicycle stock: %w",
			err,
		)
	}

	// Update booking status
	_, err = tx.Exec(
		ctx,
		`
		UPDATE bookings
		SET
			booking_status = 'cancelled',
			updated_at = CURRENT_TIMESTAMP
		WHERE booking_id = $1
		`,
		bookingID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to cancel booking: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"failed to commit cancellation: %w",
			err,
		)
	}

	return nil
}
