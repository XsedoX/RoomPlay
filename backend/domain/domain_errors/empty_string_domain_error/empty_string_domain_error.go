package empty_string_domain_error

import (
	"fmt"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
)

type EmptyStringDomainError struct {
	validation_domain_error.ValidationDomainError
}

func (e *EmptyStringDomainError) Error() string {
	return fmt.Sprintf("EmptyStringDomainError - Code: %s | Comment: %s", e.Code, e.Description)
}

func NewEmptyStringDomainError(code, fieldName string) *EmptyStringDomainError {
	validationDomainError := validation_domain_error.NewValidationDomainError(
		code+".EmptyString",
		fmt.Sprintf("The field '%s' cannot be an empty string.", fieldName),
	)
	return &EmptyStringDomainError{
		ValidationDomainError: *validationDomainError,
	}
}
