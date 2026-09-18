package service

import (
	"testing"
	"time"
)

func TestNagerHolidayService_IsHoliday(t *testing.T) {
	service := NewNagerHolidayService()

	date := time.Date(
		2026,
		time.August,
		17,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	isHoliday, err := service.IsHoliday(date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !isHoliday {
		t.Fatal("expected August 17, 2026 to be a national holiday")
	}
}
