package structure

type EmailPayload struct {
	To       string            `json:"to" validate:"required,email"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data" validate:"required"`
}
