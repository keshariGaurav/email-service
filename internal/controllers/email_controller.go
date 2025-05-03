package controllers

import (
	"context"
	"email-service/config"
	"email-service/internal/producer"
	"email-service/internal/rabbitmq"
	"email-service/structure"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func SendWelcomeEmail(c *fiber.Ctx) error {
	// Extract email data from request
	cfg := config.LoadEnv()
	rabbitConn, err := rabbitmq.NewConnection(cfg.AmqpURL)
	if err != nil {
		log.Fatal("Failed to establish RabbitMQ connection:", err)
	}
	defer rabbitConn.Close()

	prod, err := producer.NewProducer(rabbitConn.Channel, "email_queue", true)
	if err != nil {
		log.Println("Failed to create producer:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initialize producer"})
	}

	var emailData structure.EmailPayload

	if err := c.BodyParser(&emailData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if err := validate.Struct(emailData); err != nil {
		log.Println("Validation failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}


	// Publish the email data to RabbitMQ
	emailData.Template = "welcome"
	ctx := context.Background()
	err = prod.Publish(ctx, emailData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}