package time_before_now_domain_error

import (
	"fmt"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
)

type TimeBeforeNowDomainError struct {
	validation_domain_error.ValidationDomainError
}

func (e *TimeBeforeNowDomainError) Unwrap() error {
	return &e.ValidationDomainError
}

func (e *TimeBeforeNowDomainError) Error() string {
	return fmt.Sprintf("TimeBeforeNowDomainError - Code: %s | Comment: %s", e.Code, e.Description)
}

func NewTimeBeforeNowDomainError(code, fieldName string) *TimeBeforeNowDomainError {
	validationDomainError := validation_domain_error.NewValidationDomainError(
		code+".TimeBeforeNow",
		fmt.Sprintf("The field '%s' must be a time in the future.", fieldName),
	)
	return &TimeBeforeNowDomainError{
		ValidationDomainError: *validationDomainError,
	}
}
