package custom_error

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
)

type CustomError struct {
	Code      string
	ErrorType custom_error_type.Type
	Err       error
	Title     string
}

func NewCustomError(code, title string, err error, errorType custom_error_type.Type) *CustomError {
	if errors.Is(err, sql.ErrNoRows) {
		return &CustomError{Code: code, ErrorType: custom_error_type.NotFound, Err: err, Title: title}
	}
	return &CustomError{Code: code, ErrorType: errorType, Err: err, Title: title}
}

func (e CustomError) Error() string {
	var errStr string
	if e.Err != nil {
		errStr = e.Err.Error()
	}
	return fmt.Sprintf("code: %s, type: %s, error: %s, title:%s", e.Code, e.ErrorType.String(), errStr, e.Title)
}
