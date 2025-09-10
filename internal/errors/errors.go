package errors

type ErrorType string

const (
	Unauthorized ErrorType = "UNAUTHORIZED"
	Validation   ErrorType = "VALIDATION"
	NotFound     ErrorType = "NOT_FOUND"
	Internal     ErrorType = "INTERNAL"
	BadRequest   ErrorType = "BAD_REQUEST"
)

var ErrorMessages = map[ErrorType]string{
	Unauthorized: "Brak autoryzacji.",
	Validation:   "Nieprawidłowe dane.",
	NotFound:     "Zasób nie został znaleziony.",
	Internal:     "Wewnętrzny błąd serwera.",
	BadRequest:   "Błędne żądanie.",
}

type AppError struct {
	Code    string
	Type    ErrorType
	Message string
	Err     error
	Meta    map[string]any
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return ErrorMessages[e.Type]
}

var (
	ErrInvalid2FACode     = &AppError{Code: "INVALID_2FA", Type: Validation, Message: "Invalid 2FA code"}
	ErrInvalidRequest     = &AppError{Code: "INVALID_REQUEST", Type: Validation, Message: "Invalid request data"}
	ErrEmailExists        = &AppError{Code: "EMAIL_EXISTS", Type: Validation, Message: "Email already registered"}
	ErrUsernameExists     = &AppError{Code: "USERNAME_EXISTS", Type: Validation, Message: "Username already exist"}
	ErrPasswordTooShort   = &AppError{Code: "PASSWORD_TOO_SHORT", Type: Validation, Message: "Password must be at least 8 characters"}
	ErrInternal           = &AppError{Code: "SERVER_ERROR", Type: Internal, Message: "Internal server error"}
	ErrInvalidJSON        = &AppError{Code: "INVALID_JSON", Type: BadRequest, Message: "Invalid JSON in request body"}
	ErrValidationFailed   = &AppError{Code: "VALIDATION_FAILED", Type: Validation, Message: "Request validation failed"}
	ErrCSRFInvalid        = &AppError{Code: "CSRF_INVALID", Type: Unauthorized, Message: "CSRF token invalid or missing"}
	ErrTooManyRequests    = &AppError{Code: "TOO_MANY_REQUESTS", Type: BadRequest, Message: "Too many requests"}
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Type: Unauthorized, Message: "Incorrect login data"}
	ErrUserNotFound       = &AppError{Code: "USER_NOT_FOUND", Type: Unauthorized, Message: "User not found"}
	ErrUnauthorized       = &AppError{Code: "UNAUTHORIZED", Type: Unauthorized, Message: "Unauthorized access"}
)
