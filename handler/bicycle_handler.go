package handler

import (
	"net/http"
	"strconv"

	"bicycle-rent-api/entity"
	"bicycle-rent-api/usecase"

	"github.com/labstack/echo/v4"
)

type BicycleHandler struct {
	bicycleUsecase *usecase.BicycleUsecase
}

func NewBicycleHandler(bicycleUsecase *usecase.BicycleUsecase) *BicycleHandler {
	return &BicycleHandler{
		bicycleUsecase: bicycleUsecase,
	}
}

// Create godoc
// @Summary Create bicycle
// @Tags Bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.Bicycle true "Bicycle"
// @Success 201 {object} entity.Bicycle
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/bicycles [post]
func (h *BicycleHandler) Create(c echo.Context) error {
	var bicycle entity.Bicycle

	if err := c.Bind(&bicycle); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	createdBicycle, err := h.bicycleUsecase.Create(bicycle)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "bicycle created successfully",
		"data":    createdBicycle,
	})
}

// GetAll godoc
// @Summary Get all bicycles
// @Tags Bicycles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/bicycles [get]
func (h *BicycleHandler) GetAll(c echo.Context) error {
	bicycles, err := h.bicycleUsecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bicycles,
	})
}

// GetByID godoc
// @Summary Get bicycle by ID
// @Tags Bicycles
// @Produce json
// @Param id path int true "Bicycle ID"
// @Success 200 {object} entity.Bicycle
// @Failure 404 {object} map[string]interface{}
// @Router /api/bicycles/{id} [get]
func (h *BicycleHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid bicycle id",
		})
	}

	bicycle, err := h.bicycleUsecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bicycle,
	})
}

// Update godoc
// @Summary Update bicycle
// @Tags Bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Bicycle ID"
// @Param request body entity.Bicycle true "Bicycle"
// @Success 200 {object} entity.Bicycle
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/bicycles/{id} [put]
func (h *BicycleHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid bicycle id",
		})
	}

	var bicycle entity.Bicycle

	if err := c.Bind(&bicycle); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid request body",
		})
	}

	bicycle.BicycleID = id

	updatedBicycle, err := h.bicycleUsecase.Update(bicycle)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "bicycle updated successfully",
		"data":    updatedBicycle,
	})
}

// Delete godoc
// @Summary Delete bicycle
// @Tags Bicycles
// @Security BearerAuth
// @Param id path int true "Bicycle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/bicycles/{id} [delete]
func (h *BicycleHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "invalid bicycle id",
		})
	}

	if err := h.bicycleUsecase.Delete(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "bicycle deleted successfully",
	})
}
