package routes

import (
	"email-service/internal/controllers"
	"email-service/internal/producer"

	"github.com/gofiber/fiber/v2"
)

func EmailRoutes(app *fiber.App, producer *producer.Producer) {
	app.Post("/send-welcome-email", controllers.NewEmailController(producer).SendWelcomeEmail)
}
