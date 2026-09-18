package entity

import "time"

type Booking struct {
	BookingID      int64     `json:"booking_id"`
	UserID         int64     `json:"user_id"`
	BicycleID      int64     `json:"bicycle_id"`
	VoucherID      *int64    `json:"voucher_id,omitempty"`
	BookingDate    time.Time `json:"booking_date"`
	BookingStatus  string    `json:"booking_status"`
	RentalType     string    `json:"rental_type"`
	Duration       int       `json:"duration"`
	TimeDepart     time.Time `json:"time_depart"`
	TimeArrive     time.Time `json:"time_arrive"`
	SubtotalAmount float64   `json:"subtotal_amount"`
	DiscountAmount float64   `json:"discount_amount"`
	TotalAmount    float64   `json:"total_amount"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
