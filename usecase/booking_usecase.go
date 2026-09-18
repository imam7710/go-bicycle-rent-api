package usecase

import (
	"errors"
	"strings"
	"time"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/repository"
	"bicycle-rent-api/service"
)

type BookingResult struct {
	Booking      entity.Booking      `json:"booking"`
	Weather      service.WeatherInfo `json:"weather"`
	ConvertedUSD float64             `json:"converted_usd"`
}

type BookingUsecase struct {
	bookingRepository repository.BookingRepository
	bicycleRepository repository.BicycleRepository
	voucherRepository repository.VoucherRepository

	holidayService  service.HolidayService
	weatherService  service.WeatherService
	currencyService service.CurrencyService
}

func NewBookingUsecase(
	bookingRepository repository.BookingRepository,
	bicycleRepository repository.BicycleRepository,
	userRepository repository.UserRepository,
	voucherRepository repository.VoucherRepository,
	holidayService service.HolidayService,
	weatherService service.WeatherService,
	currencyService service.CurrencyService,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepository: bookingRepository,
		bicycleRepository: bicycleRepository,
		voucherRepository: voucherRepository,
		holidayService:    holidayService,
		weatherService:    weatherService,
		currencyService:   currencyService,
	}
}

func (u *BookingUsecase) Create(booking entity.Booking) (BookingResult, error) {

	if booking.UserID <= 0 {
		return BookingResult{}, errors.New("invalid user")
	}

	if booking.BicycleID <= 0 {
		return BookingResult{}, errors.New("invalid bicycle")
	}

	if booking.Duration <= 0 {
		return BookingResult{}, errors.New(
			"duration must be greater than 0",
		)
	}

	booking.RentalType = strings.ToLower(
		strings.TrimSpace(booking.RentalType),
	)

	if booking.RentalType != "day" {
		return BookingResult{}, errors.New(
			"rental type must be day",
		)
	}

	// Check bicycle
	bicycle, err := u.bicycleRepository.GetByID(
		booking.BicycleID,
	)
	if err != nil {
		return BookingResult{}, err
	}

	if bicycle.BicycleStock <= 0 {
		return BookingResult{}, errors.New(
			"bicycle is out of stock",
		)
	}

	// Check booking date
	if booking.BookingDate.IsZero() {
		return BookingResult{}, errors.New(
			"booking date is required",
		)
	}

	// Check rental date
	if booking.TimeDepart.IsZero() {
		return BookingResult{}, errors.New(
			"time depart is required",
		)
	}

	// Check whether rental date is holiday
	rentalDate := booking.TimeDepart

	isHoliday, err := u.holidayService.IsHoliday(rentalDate)
	if err != nil {
		return BookingResult{}, errors.New(
			"failed to check public holiday",
		)
	}

	// Determine daily price
	dailyPrice := bicycle.PricePerDayNormal

	if rentalDate.Weekday() == time.Saturday ||
		rentalDate.Weekday() == time.Sunday ||
		isHoliday {
		dailyPrice = bicycle.PricePerDayHoliday
	}

	subtotal := dailyPrice * float64(booking.Duration)

	// Voucher
	discountAmount := float64(0)

	if booking.VoucherID != nil {
		voucher, err := u.voucherRepository.GetByID(
			*booking.VoucherID,
		)
		if err != nil {
			return BookingResult{}, errors.New(
				"voucher not found",
			)
		}

		if !voucher.IsActive {
			return BookingResult{}, errors.New(
				"voucher is not active",
			)
		}

		if voucher.ValidUntil.Before(time.Now()) {
			return BookingResult{}, errors.New(
				"voucher has expired",
			)
		}

		discountAmount =
			subtotal * voucher.DiscountPercentage / 100

		if discountAmount > voucher.MaxDiscountAmount {
			discountAmount = voucher.MaxDiscountAmount
		}

		if discountAmount > subtotal {
			discountAmount = subtotal
		}
	}

	totalAmount := subtotal - discountAmount

	// Weather information
	weather, err := u.weatherService.GetWeather("Jakarta")
	if err != nil {
		return BookingResult{}, errors.New(
			"failed to get weather information",
		)
	}

	// Currency conversion
	convertedUSD, err := u.currencyService.Convert(
		totalAmount,
		"IDR",
		"USD",
	)
	if err != nil {
		return BookingResult{}, errors.New(
			"failed to convert currency",
		)
	}

	booking.BookingStatus = "confirmed"
	booking.SubtotalAmount = subtotal
	booking.DiscountAmount = discountAmount
	booking.TotalAmount = totalAmount

	// Atomic transaction:
	// deduct balance + stock + create booking
	createdBooking, err :=
		u.bookingRepository.CreateTransaction(
			booking,
			totalAmount,
		)

	if err != nil {
		return BookingResult{}, err
	}

	return BookingResult{
		Booking:      createdBooking,
		Weather:      weather,
		ConvertedUSD: convertedUSD,
	}, nil
}

func (u *BookingUsecase) GetByID(bookingID int64) (entity.Booking, error) {

	if bookingID <= 0 {
		return entity.Booking{}, errors.New(
			"invalid booking id",
		)
	}

	return u.bookingRepository.GetByID(bookingID)
}

func (u *BookingUsecase) GetByUserID(userID int64) ([]entity.Booking, error) {

	if userID <= 0 {
		return nil, errors.New(
			"invalid user",
		)
	}

	return u.bookingRepository.GetByUserID(userID)
}

func (u *BookingUsecase) UpdateStatus(bookingID int64, status string) error {

	if bookingID <= 0 {
		return errors.New("invalid booking id")
	}

	status = strings.ToLower(strings.TrimSpace(status))

	switch status {
	case "confirmed", "completed":
		return u.bookingRepository.UpdateStatus(
			bookingID,
			status,
		)

	case "cancelled":
		return u.bookingRepository.Cancel(bookingID)

	default:
		return errors.New("invalid booking status")
	}
}

func (u *BookingUsecase) GetReport() (repository.BookingReport, error) {
	return u.bookingRepository.GetReport()
}
