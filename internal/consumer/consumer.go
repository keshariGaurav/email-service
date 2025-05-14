package consumer

import (
	"email-service/config"
	"email-service/internal/email"
	customErrors "email-service/internal/errors"
	"email-service/internal/rabbitmq"
	"email-service/internal/retry"
	"email-service/structure"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeMessages(cfg config.Config) {
	connection, err := rabbitmq.NewConnection(cfg.AmqpURL)
	if err != nil {
		log.Fatalf("Failed to establish RabbitMQ connection: %v", err)
	}
	defer connection.Close()

	msgs, err := connection.Channel.Consume(
		connection.GetQueueName(),
		"",    // consumer
		false, // auto-ack set to false for manual acknowledgment
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Println("Consumer started, waiting for messages...")

	for msg := range msgs {
		go processMessage(msg, cfg)
	}
}

func processMessage(msg amqp.Delivery, cfg config.Config) {
	var payload structure.EmailPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		// Reject malformed messages without requeue
		msg.Reject(false)
		return
	}

	err := retry.RetryWithBackoff(func() error {
		if err := email.SendEmail(cfg, payload.To, payload.Subject, payload.Template, payload.Data); err != nil {
			if emailErr, ok := err.(*customErrors.EmailError); ok {
				log.Printf("Email error: [%s] %s", emailErr.Operation, emailErr.Message)
			}
			return err
		}
		return nil
	}, 3)

	if err != nil {
		log.Printf("Failed to process message after retries: %v", err)
		// Message processing failed after retries - reject and discard
		msg.Reject(false)
		return
	}

	// Acknowledge successful processing
	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}
