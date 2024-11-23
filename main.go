package main

import (
	"log"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Initialize routes
	routes.SetupRoutes(app)

	// Load configuration
	cfg := config.LoadConfig()

	// Start server
	log.Fatal(app.Listen(cfg.ServerAddress))
}
