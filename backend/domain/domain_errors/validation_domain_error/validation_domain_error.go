package validation_domain_error

import (
	"fmt"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
)

type ValidationDomainError struct {
	domain_errors.DomainError
	Title string
}

func (e *ValidationDomainError) Unwrap() error {
	return e.DomainError
}

func (e *ValidationDomainError) Error() string {
	return fmt.Sprintf("ValidationDomainError - Code: %s | Comment: %s", e.Code, e.Description)
}

func NewValidationDomainError(code, description string) *ValidationDomainError {
	return &ValidationDomainError{
		DomainError: domain_errors.DomainError{Code: code, Description: description},
		Title:       "Validation error occurred.",
	}
}
