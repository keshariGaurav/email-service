package routes

import (
	"email-service/config"
	"email-service/internal/controllers"
	"email-service/internal/middleware"
	"email-service/internal/producer"

	"github.com/gofiber/fiber/v2"
)

func EmailRoutes(app *fiber.App, producer *producer.Producer, cfg config.Config) {
	// Create an email group with authentication middleware
	emailGroup := app.Group("/email")
	
	// Apply JWT authentication
	emailGroup.Use(middleware.AuthMiddleware(cfg))

	// Protected routes with role-based access
	// Only users with 'admin' or 'email_service' roles can send welcome emails
	emailGroup.Post("/welcome", 
		middleware.RequireRole([]string{"admin", "email_service"}),
		controllers.NewEmailController(producer).SendWelcomeEmail,
	)
}
