package service

import "time"

type HolidayService interface {
	IsHoliday(date time.Time) (bool, error)
}
