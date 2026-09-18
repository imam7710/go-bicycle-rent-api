package handler

import (
	"net/http"

	"bicycle-rent-api/usecase"

	"github.com/labstack/echo/v4"
)

type TopupHandler struct {
	topupUsecase *usecase.TopupUsecase
}

func NewTopupHandler(topupUsecase *usecase.TopupUsecase) *TopupHandler {
	return &TopupHandler{
		topupUsecase: topupUsecase,
	}
}

type topupRequest struct {
	Amount float64 `json:"amount"`
}

// Create godoc
// @Summary Top up balance
// @Description Add dummy balance to the authenticated user's account
// @Tags Topups
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body topupRequest true "Top up request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/topups [post]
func (h *TopupHandler) Create(c echo.Context) error {
	var req topupRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "invalid user",
		})
	}

	topup, err := h.topupUsecase.Create(userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "top up successful",
		"data": map[string]interface{}{
			"topup_id":   topup.TopupID,
			"user_id":    topup.UserID,
			"amount":     topup.Amount,
			"created_at": topup.CreatedAt,
		},
	})
}

// GetHistory godoc
// @Summary Get top up history
// @Description Get top up history of the authenticated user
// @Tags Topups
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/topups [get]
func (h *TopupHandler) GetHistory(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "invalid user",
		})
	}

	topups, err := h.topupUsecase.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    topups,
	})
}
