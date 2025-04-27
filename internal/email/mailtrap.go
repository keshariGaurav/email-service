package email

import (
	"fmt"
	"net/smtp"
	"html/template"
	"bytes"
	"log"
	"email-service/config"
)

func sendWithMailtrap(cfg config.Config, to, subject, templateName string, data map[string]string) error {
	from := cfg.SenderEmail
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	htmlBody, err := RenderTemplate(templateName, data)
	if err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		to, subject, htmlBody))

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

func RenderTemplate(name string, data map[string]string) (string, error) {
	tmplPath := fmt.Sprintf("templates/%s.html", name)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Template parsing error:", err)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Println("Template execution error:", err)
		return "", err
	}

	return buf.String(), nil
}
