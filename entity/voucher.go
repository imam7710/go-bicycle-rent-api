package entity

import "time"

type Voucher struct {
	VoucherID          int64     `json:"voucher_id"`
	VoucherCode        string    `json:"voucher_code"`
	DiscountPercentage float64   `json:"discount_percentage"`
	MaxDiscountAmount  float64   `json:"max_discount_amount"`
	ValidUntil         time.Time `json:"valid_until"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
}
