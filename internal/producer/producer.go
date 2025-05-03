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
	exchange  string
	queueName string
	durable   bool
}

// NewProducer initializes a new Producer, declares the queue once, and optionally binds to an exchange.
func NewProducer(ch *amqp.Channel, exchange, queueName string, durable bool) (*Producer, error) {
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

	// Bind the queue to an exchange if specified
	if exchange != "" {
		err = ch.QueueBind(
			queueName, // queue name
			queueName, // routing key
			exchange,  // exchange
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("queue bind failed: %w", err)
		}
	}

	log.Printf("✅ Producer initialized for queue [%s] and exchange [%s]", queueName, exchange)

	return &Producer{
		channel:   ch,
		exchange:  exchange,
		queueName: queueName,
		durable:   durable,
	}, nil
}

// Publish sends a JSON-encoded message to the queue or exchange.
func (p *Producer) Publish(ctx context.Context, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,  // exchange name ("" for default)
		p.queueName, // routing key
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

	log.Printf("📤 Message published to exchange [%s], routing key [%s]", p.exchange, p.queueName)
	return nil
}
