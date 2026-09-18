// @title Bicycle Rent API
// @version 1.0
// @description REST API for Bicycle Rental System
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"log"

	"bicycle-rent-api/config"
	"bicycle-rent-api/handler"
	"bicycle-rent-api/middleware"
	"bicycle-rent-api/repository"
	"bicycle-rent-api/service"
	"bicycle-rent-api/usecase"

	_ "bicycle-rent-api/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// =========================
	// Database Connection
	// =========================
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// =========================
	// Repository Layer
	// =========================
	userRepository := repository.NewUserRepositoryPostgres(db)
	topupRepository := repository.NewTopupRepositoryPostgres(db)
	bicycleRepository := repository.NewBicycleRepositoryPostgres(db)
	bookingRepository := repository.NewBookingRepositoryPostgres(db)
	voucherRepository := repository.NewVoucherRepositoryPostgres(db)

	// =========================
	// Service Layer
	// =========================
	holidayService := service.NewNagerHolidayService()
	weatherService := service.NewOpenMeteoWeatherService()
	currencyService := service.NewFrankfurterCurrencyService()

	// =========================
	// Usecase Layer
	// =========================
	userUsecase := usecase.NewUserUsecase(userRepository)
	topupUsecase := usecase.NewTopupUsecase(topupRepository)
	bicycleUsecase := usecase.NewBicycleUsecase(bicycleRepository)
	bookingUsecase := usecase.NewBookingUsecase(
		bookingRepository,
		bicycleRepository,
		userRepository,
		voucherRepository,
		holidayService,
		weatherService,
		currencyService,
	)

	// =========================
	// Handler Layer
	// =========================
	userHandler := handler.NewUserHandler(userUsecase)
	topupHandler := handler.NewTopupHandler(topupUsecase)
	bicycleHandler := handler.NewBicycleHandler(bicycleUsecase)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)

	// =========================
	// Echo Instance
	// =========================
	e := echo.New()

	// =========================
	// Routes
	// =========================
	// User
	e.POST("/api/users/register", userHandler.Register)
	e.POST("/api/users/login", userHandler.Login)
	e.GET("/api/users/me", userHandler.GetMe, middleware.JWTMiddleware)

	// Top Up
	e.POST("/api/topups", topupHandler.Create, middleware.JWTMiddleware)
	e.GET("/api/topups", topupHandler.GetHistory, middleware.JWTMiddleware)

	// Bicycle
	e.POST("/api/bicycles", bicycleHandler.Create, middleware.JWTMiddleware)
	e.GET("/api/bicycles", bicycleHandler.GetAll)
	e.GET("/api/bicycles/:id", bicycleHandler.GetByID)
	e.PUT("/api/bicycles/:id", bicycleHandler.Update, middleware.JWTMiddleware)
	e.DELETE("/api/bicycles/:id", bicycleHandler.Delete, middleware.JWTMiddleware)

	// Booking
	e.POST("/api/bookings", bookingHandler.Create, middleware.JWTMiddleware)
	e.GET("/api/bookings", bookingHandler.GetMyBookings, middleware.JWTMiddleware)
	e.GET("/api/bookings/report", bookingHandler.GetReport, middleware.JWTMiddleware)
	e.GET("/api/bookings/:id", bookingHandler.GetByID, middleware.JWTMiddleware)
	e.PATCH("/api/bookings/:id/status", bookingHandler.UpdateStatus, middleware.JWTMiddleware)
	e.GET("/api/bookings/report", bookingHandler.GetReport, middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	// Swagger Docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// =========================
	// Server Start
	// =========================
	log.Println("Server running on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
