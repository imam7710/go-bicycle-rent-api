package handler

import (
	"net/http"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/helper"
	"bicycle-rent-api/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

type registerRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary Register user
// @Description Register a new customer account
// @Tags Users
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/users/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req registerRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	user := entity.User{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	createdUser, err := h.userUsecase.Register(user)
	if err != nil {
		status := http.StatusBadRequest

		if err.Error() == "email already registered" {
			status = http.StatusConflict
		}

		return c.JSON(status, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "user registered successfully",
		"data": map[string]interface{}{
			"user_id":   createdUser.UserID,
			"user_name": createdUser.UserName,
			"email":     createdUser.Email,
			"role":      createdUser.Role,
			"balance":   createdUser.Balance,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Login using email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var req loginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	user, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	token, err := helper.GenerateToken(user.UserID, user.Role)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "login successful",
		"data": map[string]interface{}{
			"user_id": user.UserID,
			"email":   user.Email,
			"role":    user.Role,
			"token":   token,
		},
	})
}

// GetMe godoc
// @Summary Get current user
// @Description Get current authenticated user's profile and balance
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "unauthorized",
		})
	}

	user, err := h.userUsecase.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user_id":   user.UserID,
			"user_name": user.UserName,
			"email":     user.Email,
			"role":      user.Role,
			"balance":   user.Balance,
		},
	})
}
