package shared

import "github.com/google/uuid"

// GenerateCSRFToken zwraca losowy token CSRF
func GenerateUuid() string {
	return uuid.NewString()
}
