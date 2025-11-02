package applicationErrors

import "fmt"

type ApplicationError struct {
	Code      string
	ErrorType Type
	Err       error
}

func NewApplicationError(code string, err error, errorType Type) *ApplicationError {
	return &ApplicationError{Code: code, ErrorType: errorType, Err: err}
}

func (e ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, type: %s, error: %s", e.Code, e.ErrorType.String(), e.Err.Error())
	}
	return fmt.Sprintf("code: %s, type: %s", e.Code, e.ErrorType.String())
}
