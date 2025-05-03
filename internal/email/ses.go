package email

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	customErrors "email-service/internal/errors"
	"email-service/internal/template"
)

func sendWithSES(cfg SESConfig, to, subject, templateName string, data map[string]string) error {
	// Render HTML email body from template
	htmlBody, err := template.RenderTemplate(templateName, data)
	if err != nil {
		return customErrors.NewEmailError("SES", "template rendering failed", err)
	}

	// Load AWS config with static credentials
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AwsRegion),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AwsKey, cfg.AwsSecret, ""),
		),
	)
	if err != nil {
		log.Println("Error loading AWS config:", err)
		return customErrors.NewEmailError("SES", "AWS configuration failed", err)
	}

	// Initialize SES client
	client := ses.NewFromConfig(awsCfg)

	// Construct email input
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(cfg.SenderEmail),
	}

	// Send the email
	_, err = client.SendEmail(context.TODO(), input)
	if err != nil {
		log.Println("Error sending email through SES:", err)
		return customErrors.NewEmailError("SES", "email sending failed", err)
	}

	return nil
}
