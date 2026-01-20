package domainErrors

import "fmt"

type ValidationDomainError struct {
	DomainError
	Title string
}

func (e *ValidationDomainError) Error() string {
	return fmt.Sprintf("ValidationDomainError - Code: %s | Comment: %s", e.Code, e.Description)
}

func NewValidationDomainError(code, description string) *ValidationDomainError {
	return &ValidationDomainError{
		DomainError: DomainError{Code: code, Description: description},
		Title:       "Validation error occurred.",
	}
}
