package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel   *amqp.Channel
	queueName string
	durable   bool
}

// NewProducer initializes a new Producer and declares the queue once.
func NewProducer(ch *amqp.Channel, queueName string, durable bool) (*Producer, error) {
	// Declare the queue without exchange binding
	_, err := ch.QueueDeclare(
		queueName,
		durable, // durable
		false,   // auto-delete
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	log.Printf("✅ Producer initialized for queue [%s]", queueName)

	return &Producer{
		channel:   ch,
		queueName: queueName,
		durable:   durable,
	}, nil
}

// Publish sends a JSON-encoded message to the queue.
func (p *Producer) Publish(ctx context.Context, payload interface{}) error {
	body, err := json.Marshal(payload)
	if (err != nil) {
		return fmt.Errorf("marshal failed: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",          // exchange (empty for direct queue publishing)
		p.queueName, // routing key (queue name)
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Printf("📤 Message published to queue [%s]", p.queueName)
	return nil
}
