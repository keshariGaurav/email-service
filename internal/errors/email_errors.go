package errors

import "fmt"

// EmailError represents any error that occurs during email operations
type EmailError struct {
	Operation string
	Message   string
	Err       error
}

func (e *EmailError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s failed: %s - %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("%s failed: %s", e.Operation, e.Message)
}

// NewEmailError creates a new EmailError
func NewEmailError(operation, message string, err error) *EmailError {
	return &EmailError{
		Operation: operation,
		Message:   message,
		Err:       err,
	}
}