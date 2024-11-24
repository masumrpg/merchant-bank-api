package routes

import (
	"merchant-bank-api/app/handlers"
	"merchant-bank-api/app/middleware"
	"merchant-bank-api/app/repository"
	"merchant-bank-api/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Repo layer
	userRepo := repository.NewUserRepository()
	activityRepo := repository.NewHistoryRepository()
	transactionRepo := repository.NewTransactionRepository()

	// Service layer
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, activityRepo)
	paymentService := services.NewPaymentService(transactionRepo, userRepo)
	activityService := services.NewHistoryService(activityRepo)

	// Handler layer
	jwtMiddlewareRepo := middleware.NewJWTHandler(userRepo, activityRepo)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	activityHandler := handlers.NewHistoryHandler(activityService)

	// Public routes
	app.Post("/api/login", authHandler.Login)
	app.Post("/api/register", authHandler.Register)

	// Protected routes
	app.Use("/api/protected", middleware.JWTMiddleware(jwtMiddlewareRepo))
	app.Post("/api/protected/logout", authHandler.Logout)
	// Payment routes
	app.Post("/api/protected/payments", paymentHandler.ProcessPayment)
	app.Get("/api/protected/payments", paymentHandler.GetAllPayment)
	app.Get("/api/protected/payments/:id", paymentHandler.GetPaymentById)
	// Activities routes
	app.Get("/api/protected/activities", activityHandler.GetAllActivities)
	app.Get("/api/protected/activities/:username", activityHandler.GetActivitiesByUsername)
	// User routes
	app.Get("/api/protected/users", userHandler.GetAllUser)
	app.Get("/api/protected/users/:id", userHandler.GetUserById)

}
