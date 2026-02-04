package domain_errors

import "fmt"

type DomainError struct {
	Code        string
	Description string
}

func (d DomainError) Error() string {
	return fmt.Sprintf("DomainError: Code: %s | Comment: %s", d.Code, d.Description)
}
