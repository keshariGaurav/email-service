package rabbitmq

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	mu         sync.Mutex
	Conn       *amqp.Connection
	Channel    *amqp.Channel
	amqpURL    string
	notifyConn chan *amqp.Error
	notifyChan chan *amqp.Error
	ctx    		 context.Context
}

var (
	instance *Connection
	once     sync.Once
)

// NewConnection ensures singleton pattern with reconnect logic
func NewConnection(amqpURL string) (*Connection, error) {
	var err error
	once.Do(func() {
		conn := &Connection{amqpURL: amqpURL}
		err = conn.connect()
		if err == nil {
			instance = conn
			go instance.reconnectOnFailure()
		}
	})
	return instance, err
}

func (c *Connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.amqpURL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	c.Conn = conn
	c.Channel = ch
	c.notifyConn = conn.NotifyClose(make(chan *amqp.Error))
	c.notifyChan = ch.NotifyClose(make(chan *amqp.Error))

	log.Println("✅ RabbitMQ connected and channel opened.")
	return nil
}

// reconnectOnFailure watches for closures and reconnects automatically
func (c *Connection) reconnectOnFailure() {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("🛑 Stopping reconnect goroutine.")
			return
		case err := <-c.notifyConn:
			log.Printf("🚨 RabbitMQ connection closed: %v. Reconnecting...", err)
			c.reconnect()
		case err := <-c.notifyChan:
			log.Printf("🚨 RabbitMQ channel closed: %v. Reconnecting...", err)
			c.reconnect()
		}
	}
}

func (c *Connection) reconnect() {
	wait := time.Second
	for {
		err := c.connect()
		if err == nil {
			return
		}
		log.Printf("Reconnection failed: %v. Retrying in %v...", err, wait)
		time.Sleep(wait)
		wait *= 2
		if wait > 30*time.Second {
			wait = 30 * time.Second
		}
}
}

func (c *Connection) Close() {
	if c.Channel != nil {
		_ = c.Channel.Close()
	}
	if c.Conn != nil {
		_ = c.Conn.Close()
	}
	instance = nil
}

func (c *Connection) Publish(exchange, routingKey string, body []byte) error {
	return c.Channel.Publish(
		exchange,    // exchange
		routingKey,  // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *Connection) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(
		queue, // queue
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
