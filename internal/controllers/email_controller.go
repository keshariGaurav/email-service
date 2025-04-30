package controllers

import (
	"email-service/config"
	"email-service/internal/producer"
	"email-service/internal/rabbitmq"
	"email-service/structure"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SendWelcomeEmail(c *fiber.Ctx) error {
	// Extract email data from request
	cfg := config.LoadEnv()
	rabbitConn, err := rabbitmq.NewConnection(cfg.AmqpURL)
	if err != nil {
		log.Fatal("Failed to establish RabbitMQ connection:", err)
	}
	defer rabbitConn.Close()

	prod := producer.NewProducer(rabbitConn.Channel, cfg.QueueName)

	var emailData structure.EmailPayload

	if err := c.BodyParser(&emailData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Publish the email data to RabbitMQ
	err = prod.Publish(emailData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}