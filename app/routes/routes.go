package routes

import (
	"merchant-bank-api/app/handlers"
	"merchant-bank-api/app/middleware"
	"merchant-bank-api/app/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	repo := repository.NewJSONRepository()

	authHandler := handlers.NewAuthHandler(repo)
	paymentHandler := handlers.NewPaymentHandler(repo)
	customerHandler := handlers.NewCustomerHandler(repo)

	// Public routes
	app.Post("/api/login", authHandler.Login)
	app.Post("/api/register", authHandler.Register)

	// Protected routes
	app.Use("/api/protected", middleware.JWTMiddleware(repo))
	app.Post("/api/protected/logout", authHandler.Logout)
	// Payment routes
	app.Post("/api/protected/payments", paymentHandler.ProcessPayment)
	app.Get("/api/protected/payments", paymentHandler.GetAllPayment)
	app.Get("/api/protected/payments/:id", paymentHandler.GetPaymentById)
	// Activities routes
	app.Get("/api/protected/activities", authHandler.GetAllActivities)
	app.Get("/api/protected/activities/:username", authHandler.GetActivityByUsername)
	// Customer routes
	app.Get("/api/protected/customers", customerHandler.GetAllCustomer)
	app.Get("/api/protected/customers/:id", customerHandler.GetCustomerById)

}
