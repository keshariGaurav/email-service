package main

import (
	"email-service/config"
	"email-service/internal/consumer"
	"email-service/internal/producer"
	"email-service/internal/rabbitmq"
	"email-service/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

var (
	rabbitConn *rabbitmq.Connection
	emailProducer *producer.Producer
)

func main() {
	cfg := config.LoadEnv()
	
	// Initialize RabbitMQ connection
	var err error
	rabbitConn, err = rabbitmq.NewConnection(cfg.AmqpURL)
	if err != nil {
		log.Fatal("Failed to establish RabbitMQ connection:", err)
	}
	defer rabbitConn.Close()

	// Initialize producer
	emailProducer, err = producer.NewProducer(rabbitConn.Channel, "email_queue", true)
	if err != nil {
		log.Fatal("Failed to create producer:", err)
	}

	// Start consumer in a goroutine
	go consumer.ConsumeMessages(cfg)

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes with producer
	routes.EmailRoutes(app, emailProducer)

	// Start server
	log.Fatal(app.Listen(":6000"))
}
