package server

import (
	"email-service/config"
	"email-service/internal/consumer"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Email Service Running")
	})

	// Call the consumer to start RabbitMQ listening
	go consumer.ConsumeMessages(config.LoadEnv())
}
