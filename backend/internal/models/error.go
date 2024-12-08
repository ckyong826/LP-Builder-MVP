package models

type Error string

func (e Error) Error() string {
	return string(e)
}

type ServiceError struct {
	Message    string
	StatusCode int
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(message string, statusCode int) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: statusCode,
	}
}