package main

import (
	"email-service/config"
	"email-service/internal/consumer"

	"email-service/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadEnv()
	go consumer.ConsumeMessages(cfg)

	

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	routes.EmailRoutes(app)

	// Start server
	log.Fatal(app.Listen(":6000"))
}
