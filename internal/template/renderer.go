package template

import (
	"bytes"
	"html/template"
	"log"
	"fmt"
)

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
