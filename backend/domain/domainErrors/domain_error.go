package domainErrors

import "fmt"

type DomainErrorType int

const (
	Validation DomainErrorType = iota
)

var domainErrorTypeToString = map[DomainErrorType]string{
	Validation: "validation",
}
var stringToDomainErrorType = map[string]DomainErrorType{
	"validation": Validation,
}

func (t DomainErrorType) String() string {
	return domainErrorTypeToString[t]
}
func ParseType(s string) (DomainErrorType, bool) {
	v, ok := stringToDomainErrorType[s]
	return v, ok
}

type DomainError struct {
	code      string
	comment   string
	errorType DomainErrorType
}

func (d DomainError) Error() string {
	return fmt.Sprintf("DomainError: Code: %s | Comment: %s | Type: %s ", d.code, d.comment, d.errorType.String())
}
func NewValidationDomainError(code, comment string) *DomainError {
	return &DomainError{code: code, comment: comment, errorType: Validation}
}
