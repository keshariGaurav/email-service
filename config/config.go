package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env           string
	AmqpURL       string
	QueueName     string

	// Mailtrap
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string

	// AWS SES
	AwsRegion     string
	AwsKey        string
	AwsSecret     string
	SenderEmail   string
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, continuing")
	}

	return Config{
		Env:         os.Getenv("ENV"),
		AmqpURL:     os.Getenv("AMQP_URL"),
		QueueName:   os.Getenv("QUEUE_NAME"),

		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASS"),

		AwsRegion:   os.Getenv("AWS_REGION"),
		AwsKey:      os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecret:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SenderEmail: os.Getenv("SENDER_EMAIL"),
	}
}
