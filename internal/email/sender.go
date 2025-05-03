package email

import (
	"log"
	"os"

	"email-service/config"
	customErrors "email-service/internal/errors"
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

// validateConfig checks if all required configuration is present
func validateConfig(cfg config.Config, env string) error {
	switch env {
	case "development":
		if cfg.SMTPHost == "" {
			return customErrors.NewEmailError("Config", "missing SMTP_HOST", nil)
		}
		if cfg.SMTPPort == "" {
			return customErrors.NewEmailError("Config", "missing SMTP_PORT", nil)
		}
		if cfg.SMTPUser == "" {
			return customErrors.NewEmailError("Config", "missing SMTP_USER", nil)
		}
		if cfg.SMTPPass == "" {
			return customErrors.NewEmailError("Config", "missing SMTP_PASS", nil)
		}
	case "production":
		if awsConfig.AwsRegion == "" {
			return customErrors.NewEmailError("Config", "missing AWS_REGION", nil)
		}
		if awsConfig.AwsKey == "" {
			return customErrors.NewEmailError("Config", "missing AWS_ACCESS_KEY_ID", nil)
		}
		if awsConfig.AwsSecret == "" {
			return customErrors.NewEmailError("Config", "missing AWS_SECRET_ACCESS_KEY", nil)
		}
		if awsConfig.SenderEmail == "" {
			return customErrors.NewEmailError("Config", "missing SENDER_EMAIL", nil)
		}
	}
	return nil
}

// validateEmailParams checks if all required email parameters are present
func validateEmailParams(to, subject, templateName string, data map[string]string) error {
	if to == "" {
		return customErrors.NewEmailError("Validation", "recipient email cannot be empty", nil)
	}
	if templateName == "" {
		return customErrors.NewEmailError("Validation", "template name cannot be empty", nil)
	}
	if data == nil {
		return customErrors.NewEmailError("Validation", "template data cannot be nil", nil)
	}
	return nil
}

// SendEmail decides which service (Mailtrap or SES) to use and handles errors appropriately
func SendEmail(cfg config.Config, to, subject, templateName string, data map[string]string) error {
	// Validate email parameters
	if err := validateEmailParams(to, subject, templateName, data); err != nil {
		return err
	}

	// Validate configuration
	if err := validateConfig(cfg, cfg.Env); err != nil {
		return err
	}

	var err error
	switch cfg.Env {
	case "development":
		log.Printf("Using Mailtrap for sending email to: %s", to)
		err = sendWithMailtrap(cfg, to, subject, templateName, data)
		if err != nil {
			return customErrors.NewEmailError("Mailtrap", "failed to send email", err)
		}
	case "production":
		log.Printf("Using Amazon SES for sending email to: %s", to)
		err = sendWithSES(awsConfig, to, subject, templateName, data)
		if err != nil {
			return customErrors.NewEmailError("SES", "failed to send email", err)
		}
	default:
		log.Printf("Unknown environment: %s. Defaulting to Mailtrap.", cfg.Env)
		err = sendWithMailtrap(cfg, to, subject, templateName, data)
		if err != nil {
			return customErrors.NewEmailError("Mailtrap", "failed to send email", err)
		}
	}

	log.Printf("Email sent successfully to: %s", to)
	return nil
}
