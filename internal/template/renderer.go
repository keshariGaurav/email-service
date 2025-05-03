package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	customErrors "email-service/internal/errors"
)

func RenderTemplate(name string, data map[string]string) (string, error) {
	tmplPath := fmt.Sprintf("templates/%s.html", name)

	// Check if template file exists
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		return "", customErrors.NewEmailError("Template", fmt.Sprintf("template file does not exist: %s", tmplPath), err)
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Failed to parse template %s: %v", name, err)
		return "", customErrors.NewEmailError("Template", "failed to parse template", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Failed to execute template %s: %v", name, err)
		return "", customErrors.NewEmailError("Template", "failed to execute template", err)
	}

	return buf.String(), nil
}
