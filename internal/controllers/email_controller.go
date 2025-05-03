package controllers

import (
	"context"
	customErrors "email-service/internal/errors"
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
		log.Printf("Failed to parse request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := validate.Struct(emailData); err != nil {
		log.Printf("Validation error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}

	emailData.Template = "welcome"
	ctx := context.Background()
	if err := ec.producer.Publish(ctx, emailData); err != nil {
		if emailErr, ok := err.(*customErrors.EmailError); ok {
			log.Printf("Email error: %v", emailErr)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   emailErr.Message,
				"details": emailErr.Operation,
			})
		}
		log.Printf("Unexpected error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process email request",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Email queued successfully",
	})
}