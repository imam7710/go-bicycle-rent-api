package repository

import "bicycle-rent-api/entity"

type BookingReport struct {
	TotalBookings     int64   `json:"total_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	CompletedBookings int64   `json:"completed_bookings"`
	CancelledBookings int64   `json:"cancelled_bookings"`
}

type BookingRepository interface {
	Create(booking entity.Booking) (entity.Booking, error)
	GetByID(bookingID int64) (entity.Booking, error)
	GetByUserID(userID int64) ([]entity.Booking, error)
	UpdateStatus(bookingID int64, status string) error
	Cancel(bookingID int64) error
	GetReport() (BookingReport, error)

	CreateTransaction(booking entity.Booking, totalAmount float64) (entity.Booking, error)
}
