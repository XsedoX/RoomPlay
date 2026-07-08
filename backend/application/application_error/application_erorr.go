package application_error

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
)

type ApplicationError struct {
	Code      string
	ErrorType application_error_type.ApplicationErrorType
	Err       error
	Title     string
}

func (e *ApplicationError) Unwrap() error {
	return e.Err
}

func NewApplicationError(code, title string, err error, errorType application_error_type.ApplicationErrorType) *ApplicationError {
	if errors.Is(err, sql.ErrNoRows) {
		return &ApplicationError{Code: code, ErrorType: application_error_type.NotFound, Err: err, Title: title}
	}
	return &ApplicationError{Code: code, ErrorType: errorType, Err: err, Title: title}
}

func (e ApplicationError) Error() string {
	var errStr string
	if e.Err != nil {
		errStr = e.Err.Error()
	}
	return fmt.Sprintf("code: %s, type: %s, error: %s, title:%s", e.Code, e.ErrorType.String(), errStr, e.Title)
}
