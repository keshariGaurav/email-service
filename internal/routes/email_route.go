package routes

import (
	"email-service/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func EmailRoutes(app *fiber.App) {
	app.Post("/send-welcome-email", controllers.SendWelcomeEmail)
}
