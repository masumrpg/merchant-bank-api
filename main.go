package main

import (
	"log"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupApp sets up the Fiber application (for testing and main runtime)
func SetupApp() *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
	}))
	app.Use(logger.New())

	// Health check route
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is up and running",
		})
	})

	// Initialize routes
	routes.SetupRoutes(app)

	return app
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup app
	app := SetupApp()

	// Start server
	log.Fatal(app.Listen(cfg.ServerAddress))
}
