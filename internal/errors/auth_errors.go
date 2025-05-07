package errors

// AuthError represents authentication related errors
type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func NewAuthError(code int, message string) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}