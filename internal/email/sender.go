package email

import (
	"fmt"
	"log"
	"os"

	"email-service/config"
)

// SESConfig holds AWS credentials for SES
type SESConfig struct {
	AwsRegion   string
	AwsKey      string
	AwsSecret   string
	SenderEmail string
}

var awsConfig SESConfig

// Initialize AWS SES Config
func init() {
	awsConfig = SESConfig{
		AwsRegion:   os.Getenv("AWS_REGION"),
		AwsKey:      os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecret:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SenderEmail: os.Getenv("SENDER_EMAIL"),
	}
}

// SendEmail decides which service (Mailtrap or SES) to use
func SendEmail(cfg config.Config, to, subject, templateName string, data map[string]string) error {
	log.Printf("Preparing to send email to: %s", to)
	fmt.Print("sending email")

	switch cfg.Env {
	case "development":
		log.Println("Environment: development. Using Mailtrap for sending email.")
		return sendWithMailtrap(cfg, to, subject, templateName, data)
	case "production":
		log.Println("Environment: production. Using Amazon SES for sending email.")
		return sendWithSES(awsConfig, to, subject, templateName, data)
	default:
		log.Printf("Unknown environment: %s. Defaulting to Mailtrap.", cfg.Env)
		return sendWithMailtrap(cfg, to, subject, templateName, data)
	}
}
