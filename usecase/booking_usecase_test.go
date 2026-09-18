package usecase

import (
	"testing"
	"time"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/repository"
	"bicycle-rent-api/service"
)

// =========================
// Mock Booking Repository
// =========================

type mockBookingRepository struct {
	createdBooking entity.Booking
	createError    error
}

func (m *mockBookingRepository) Create(
	booking entity.Booking,
) (entity.Booking, error) {
	return booking, nil
}

func (m *mockBookingRepository) GetByID(
	bookingID int64,
) (entity.Booking, error) {
	return entity.Booking{}, nil
}

func (m *mockBookingRepository) GetByUserID(
	userID int64,
) ([]entity.Booking, error) {
	return nil, nil
}

func (m *mockBookingRepository) UpdateStatus(
	bookingID int64,
	status string,
) error {
	return nil
}

func (m *mockBookingRepository) Cancel(
	bookingID int64,
) error {
	return nil
}

func (m *mockBookingRepository) GetReport() (
	repository.BookingReport,
	error,
) {
	return repository.BookingReport{}, nil
}

func (m *mockBookingRepository) CreateTransaction(
	booking entity.Booking,
	totalAmount float64,
) (entity.Booking, error) {

	if m.createError != nil {
		return entity.Booking{}, m.createError
	}

	booking.TotalAmount = totalAmount
	return booking, nil
}

// =========================
// Mock Bicycle Repository
// =========================

type mockBicycleRepository struct {
	bicycle entity.Bicycle
}

func (m *mockBicycleRepository) Create(
	bicycle entity.Bicycle,
) (entity.Bicycle, error) {
	return bicycle, nil
}

func (m *mockBicycleRepository) GetAll() (
	[]entity.Bicycle,
	error,
) {
	return nil, nil
}

func (m *mockBicycleRepository) GetByID(
	bicycleID int64,
) (entity.Bicycle, error) {
	return m.bicycle, nil
}

func (m *mockBicycleRepository) Update(
	bicycle entity.Bicycle,
) (entity.Bicycle, error) {
	return bicycle, nil
}

func (m *mockBicycleRepository) Delete(
	bicycleID int64,
) error {
	return nil
}

// =========================
// Mock User Repository
// =========================

type mockUserRepository struct{}

func (m *mockUserRepository) Create(
	user entity.User,
) (entity.User, error) {
	return user, nil
}

func (m *mockUserRepository) UpdateBalance(
	userID int64,
	balance float64,
) error {
	return nil
}

func (m *mockUserRepository) GetAll() (
	[]entity.User,
	error,
) {
	return nil, nil
}

func (m *mockUserRepository) GetByID(
	userID int64,
) (entity.User, error) {
	return entity.User{}, nil
}

func (m *mockUserRepository) GetByEmail(
	email string,
) (entity.User, error) {
	return entity.User{}, nil
}

func (m *mockUserRepository) Update(
	user entity.User,
) (entity.User, error) {
	return user, nil
}

func (m *mockUserRepository) Delete(
	userID int64,
) error {
	return nil
}

// =========================
// Mock Voucher Repository
// =========================

type mockVoucherRepository struct {
	voucher entity.Voucher
}

func (m *mockVoucherRepository) Create(
	voucher entity.Voucher,
) (entity.Voucher, error) {
	return voucher, nil
}

func (m *mockVoucherRepository) GetAll() (
	[]entity.Voucher,
	error,
) {
	return nil, nil
}

func (m *mockVoucherRepository) GetByID(
	voucherID int64,
) (entity.Voucher, error) {
	return m.voucher, nil
}

func (m *mockVoucherRepository) GetByCode(
	code string,
) (entity.Voucher, error) {
	return m.voucher, nil
}

func (m *mockVoucherRepository) Update(
	voucher entity.Voucher,
) (entity.Voucher, error) {
	return voucher, nil
}

func (m *mockVoucherRepository) Delete(
	voucherID int64,
) error {
	return nil
}

// =========================
// Mock Holiday Service
// =========================

type mockHolidayService struct {
	isHoliday bool
}

func (m *mockHolidayService) IsHoliday(
	date time.Time,
) (bool, error) {
	return m.isHoliday, nil
}

// =========================
// Mock Weather Service
// =========================

type mockWeatherService struct{}

func (m *mockWeatherService) GetWeather(
	city string,
) (service.WeatherInfo, error) {
	return service.WeatherInfo{
		City:        city,
		Temperature: 30,
		Description: "Clear sky",
	}, nil
}

// =========================
// Mock Currency Service
// =========================

type mockCurrencyService struct{}

func (m *mockCurrencyService) Convert(
	amount float64,
	from string,
	to string,
) (float64, error) {
	return amount / 17000, nil
}

// =========================
// Helper
// =========================

func newTestBookingUsecase(
	holiday bool,
	voucher *entity.Voucher,
) *BookingUsecase {

	return NewBookingUsecase(
		&mockBookingRepository{},
		&mockBicycleRepository{
			bicycle: entity.Bicycle{
				BicycleID:          1,
				BicycleStock:       5,
				PricePerDayNormal:  50000,
				PricePerDayHoliday: 75000,
			},
		},
		&mockUserRepository{},
		&mockVoucherRepository{
			voucher: func() entity.Voucher {
				if voucher != nil {
					return *voucher
				}
				return entity.Voucher{}
			}(),
		},
		&mockHolidayService{
			isHoliday: holiday,
		},
		&mockWeatherService{},
		&mockCurrencyService{},
	)
}

// =========================
// Test 1: Normal Booking
// =========================

func TestBookingUsecase_Create_Normal(t *testing.T) {

	uc := newTestBookingUsecase(false, nil)

	booking := entity.Booking{
		UserID:      1,
		BicycleID:   1,
		BookingDate: time.Now(),
		RentalType:  "day",
		Duration:    1,
		TimeDepart:  time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC),
		TimeArrive:  time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC),
	}

	result, err := uc.Create(booking)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Booking.TotalAmount != 50000 {
		t.Fatalf(
			"expected total 50000, got %.2f",
			result.Booking.TotalAmount,
		)
	}
}

// =========================
// Test 2: Holiday Booking
// =========================

func TestBookingUsecase_Create_Holiday(t *testing.T) {

	uc := newTestBookingUsecase(true, nil)

	booking := entity.Booking{
		UserID:      1,
		BicycleID:   1,
		BookingDate: time.Now(),
		RentalType:  "day",
		Duration:    1,
		TimeDepart:  time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC),
		TimeArrive:  time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC),
	}

	result, err := uc.Create(booking)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Booking.TotalAmount != 75000 {
		t.Fatalf(
			"expected total 75000, got %.2f",
			result.Booking.TotalAmount,
		)
	}
}

// =========================
// Test 3: Voucher
// =========================

func TestBookingUsecase_Create_WithVoucher(t *testing.T) {

	voucherID := int64(1)

	voucher := &entity.Voucher{
		VoucherID:          voucherID,
		VoucherCode:        "HEMAT10",
		DiscountPercentage: 10,
		MaxDiscountAmount:  10000,
		ValidUntil:         time.Now().Add(24 * time.Hour),
		IsActive:           true,
	}

	uc := newTestBookingUsecase(false, voucher)

	booking := entity.Booking{
		UserID:      1,
		BicycleID:   1,
		VoucherID:   &voucherID,
		BookingDate: time.Now(),
		RentalType:  "day",
		Duration:    1,
		TimeDepart:  time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC),
		TimeArrive:  time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC),
	}

	result, err := uc.Create(booking)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Booking.SubtotalAmount != 50000 {
		t.Fatalf(
			"expected subtotal 50000, got %.2f",
			result.Booking.SubtotalAmount,
		)
	}

	if result.Booking.DiscountAmount != 5000 {
		t.Fatalf(
			"expected discount 5000, got %.2f",
			result.Booking.DiscountAmount,
		)
	}

	if result.Booking.TotalAmount != 45000 {
		t.Fatalf(
			"expected total 45000, got %.2f",
			result.Booking.TotalAmount,
		)
	}
}

func (m *mockBicycleRepository) UpdateStock(
	bicycleID int64,
	stock int,
) error {
	return nil
}

func TestBookingUsecase_Create_InactiveVoucher(t *testing.T) {
	voucherID := int64(1)

	voucher := &entity.Voucher{
		VoucherID:          voucherID,
		VoucherCode:        "INACTIVE10",
		DiscountPercentage: 10,
		MaxDiscountAmount:  10000,
		ValidUntil:         time.Now().Add(24 * time.Hour),
		IsActive:           false,
	}

	uc := newTestBookingUsecase(false, voucher)

	booking := entity.Booking{
		UserID:      1,
		BicycleID:   1,
		VoucherID:   &voucherID,
		BookingDate: time.Now(),
		RentalType:  "day",
		Duration:    1,
		TimeDepart:  time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC),
		TimeArrive:  time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC),
	}

	_, err := uc.Create(booking)

	if err == nil {
		t.Fatal("expected inactive voucher error")
	}

	if err.Error() != "voucher is not active" {
		t.Fatalf(
			"expected inactive voucher error, got %v",
			err,
		)
	}
}

func TestBookingUsecase_Create_ExpiredVoucher(t *testing.T) {
	voucherID := int64(1)

	voucher := &entity.Voucher{
		VoucherID:          voucherID,
		VoucherCode:        "EXPIRED10",
		DiscountPercentage: 10,
		MaxDiscountAmount:  10000,
		ValidUntil:         time.Now().Add(-24 * time.Hour),
		IsActive:           true,
	}

	uc := newTestBookingUsecase(false, voucher)

	booking := entity.Booking{
		UserID:      1,
		BicycleID:   1,
		VoucherID:   &voucherID,
		BookingDate: time.Now(),
		RentalType:  "day",
		Duration:    1,
		TimeDepart:  time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC),
		TimeArrive:  time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC),
	}

	_, err := uc.Create(booking)

	if err == nil {
		t.Fatal("expected expired voucher error")
	}

	if err.Error() != "voucher has expired" {
		t.Fatalf(
			"expected expired voucher error, got %v",
			err,
		)
	}
}
