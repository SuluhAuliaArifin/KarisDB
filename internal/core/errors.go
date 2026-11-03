package core

import "fmt"

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
