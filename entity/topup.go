package entity

import "time"

type Topup struct {
	TopupID   int64     `json:"topup_id"`
	UserID    int64     `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
