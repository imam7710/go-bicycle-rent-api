package entity

import "time"

type Bicycle struct {
	BicycleID          int64     `json:"bicycle_id"`
	BicycleModel       string    `json:"bicycle_model"`
	BicycleType        string    `json:"bicycle_type"`
	BicycleStock       int       `json:"bicycle_stock"`
	PricePerDayNormal  float64   `json:"price_per_day_normal"`
	PricePerDayHoliday float64   `json:"price_per_day_holiday"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
