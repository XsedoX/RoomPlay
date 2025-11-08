package applicationErrors

import "fmt"

type ApplicationError struct {
	Code      string
	ErrorType Type
	Err       error
	Title     string
}

func NewApplicationError(code, title string, err error, errorType Type) *ApplicationError {
	return &ApplicationError{Code: code, ErrorType: errorType, Err: err, Title: title}
}

func (e ApplicationError) Error() string {
	var errStr string
	if e.Err != nil {
		errStr = e.Err.Error()
	}
	return fmt.Sprintf("code: %s, type: %s, error: %s, title:%s", e.Code, e.ErrorType.String(), errStr, e.Title)
}
