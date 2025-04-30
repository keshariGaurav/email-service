package email

import (
	"email-service/config"
	"email-service/internal/template"
	"fmt"
	"net/smtp"
)

func sendWithMailtrap(cfg config.Config, to, subject, templateName string, data map[string]string) error {
	from := cfg.SenderEmail
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	htmlBody, err := template.RenderTemplate(templateName, data)
	if err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		to, subject, htmlBody))

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

