package main

import (
	"email-service/config"
	"email-service/internal/consumer"
	"email-service/internal/server"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadEnv()
	go consumer.ConsumeMessages(cfg)

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	server.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
