package main

import (
	"email-service/config"
	"email-service/internal/consumer"
	"email-service/internal/rabbitmq"
	"log"

	"github.com/gofiber/fiber/v2"
)

var (
	rabbitConn    *rabbitmq.Connection

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

	// Start consumer in a goroutine
	go consumer.ConsumeMessages(cfg)

	// Initialize Fiber app
	app := fiber.New()

	// Start server
	log.Fatal(app.Listen(":6000"))
}
