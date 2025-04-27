package consumer

import (
	"encoding/json"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"email-service/config"
	"email-service/internal/email"
	"email-service/internal/retry"
)


type EmailPayload struct {
	To       string            `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

func ConsumeMessages(cfg config.Config) {
	conn, err := amqp.Dial(cfg.AmqpURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Queue declare error:", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to consume messages:", err)
	}

	for msg := range msgs {
		go func(d amqp.Delivery) {
			var payload EmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Println("Invalid payload:", err)
				return
			}

			// Retry sending email with exponential backoff
			retry.RetryWithBackoff(func() error {
				return email.SendEmail(cfg, payload.To, payload.Subject, payload.Template, payload.Data)
			}, 3)
		}(msg)
	}
}
