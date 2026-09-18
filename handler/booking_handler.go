package handler

import (
	"net/http"
	"strconv"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/usecase"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingUsecase *usecase.BookingUsecase
}

func NewBookingHandler(bookingUsecase *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{bookingUsecase: bookingUsecase}
}

// Create godoc
// @Summary Create booking
// @Description Create a bicycle rental booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body entity.Booking true "Booking request"
// @Success 201 {object} usecase.BookingResult
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/bookings [post]
func (h *BookingHandler) Create(c echo.Context) error {
	var booking entity.Booking

	if err := c.Bind(&booking); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "invalid user token",
		})
	}

	booking.UserID = userID

	result, err := h.bookingUsecase.Create(booking)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "booking created successfully",
		"data":    result,
	})
}

// GetByID godoc
// @Summary Get booking by ID
// @Tags Bookings
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} entity.Booking
// @Failure 404 {object} map[string]interface{}
// @Router /api/bookings/{id} [get]
func (h *BookingHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid booking id",
		})
	}

	booking, err := h.bookingUsecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    booking,
	})
}

// GetMyBookings godoc
// @Summary Get my bookings
// @Tags Bookings
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/bookings [get]
func (h *BookingHandler) GetMyBookings(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "invalid user token",
		})
	}

	bookings, err := h.bookingUsecase.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bookings,
	})
}

// UpdateStatus godoc
// @Summary Update booking status
// @Tags Bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Booking ID"
// @Param request body map[string]string true "Booking status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/bookings/{id}/status [patch]
func (h *BookingHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid booking id",
		})
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := h.bookingUsecase.UpdateStatus(id, request.Status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "booking status updated successfully",
	})
}

// GetReport godoc
// @Summary Get booking report
// @Description Get booking report (admin only)
// @Tags Bookings
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} repository.BookingReport
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/bookings/report [get]
func (h *BookingHandler) GetReport(c echo.Context) error {
	report, err := h.bookingUsecase.GetReport()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}
