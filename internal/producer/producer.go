package producer

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel   *amqp.Channel
	queueName string
}

// NewProducer creates a new producer instance
func NewProducer(ch *amqp.Channel, queueName string) *Producer {
	return &Producer{
		channel:   ch,
		queueName: queueName,
	}
}

// Publish sends a message to the RabbitMQ queue
func (p *Producer) Publish(payload interface{}) error {
    // Declare the queue first
    _, err := p.channel.QueueDeclare(
        p.queueName,
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
		fmt.Println("Queue declared:", p.queueName)
    if err != nil {
        log.Println("Failed to declare queue:", err)
        return err
    }

    body, err := json.Marshal(payload)
    if err != nil {
        log.Println("Failed to marshal payload:", err)
        return err
    }

    err = p.channel.Publish(
        "",          // Default exchange
        p.queueName, // Routing key (queue name)
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )

    if err != nil {
        log.Println("Failed to publish message:", err)
        return err
    }

    log.Println("Message published successfully to queue:", p.queueName)
    return nil
}
