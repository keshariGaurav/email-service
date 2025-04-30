package consumer

import (
	"email-service/config"
	"email-service/internal/email"
	"email-service/internal/rabbitmq"
	"email-service/structure"
	"email-service/internal/retry"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)



func ConsumeMessages(cfg config.Config) {
	rabbitConn, err := rabbitmq.NewConnection(cfg.AmqpURL)
	fmt.Println("RabbitMQ connection established")
	if err != nil {
		log.Fatal("Failed to establish RabbitMQ connection:", err)
	}
	defer rabbitConn.Close()

	msgs, err := rabbitConn.Channel.Consume(
		"email-queue", 
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatal("Failed to consume messages:", err)
	}
	fmt.Println("Waiting for messages...", msgs)

	for msg := range msgs {
		go func(d amqp.Delivery) {
			var payload structure.EmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("Invalid payload:", err)
				return
			}

			// Retry sending email with exponential backoff
			retry.RetryWithBackoff(func() error {
				fmt.Print("Sending email to:", payload.To)
				return email.SendEmail(cfg, payload.To, payload.Subject, payload.Template, payload.Data)
			}, 3)
		}(msg)
	}
}
