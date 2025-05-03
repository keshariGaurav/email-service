package controllers

import (
	"context"
	"email-service/internal/producer"
	"email-service/structure"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type EmailController struct {
	producer *producer.Producer
}

func NewEmailController(producer *producer.Producer) *EmailController {
	return &EmailController{
		producer: producer,
	}
}

func (ec *EmailController) SendWelcomeEmail(c *fiber.Ctx) error {
	var emailData structure.EmailPayload
	if err := c.BodyParser(&emailData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if err := validate.Struct(emailData); err != nil {
		log.Printf("Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Publish the email data to RabbitMQ
	emailData.Template = "welcome"
	ctx := context.Background()
	err := ec.producer.Publish(ctx, emailData)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Email sent successfully"})
}