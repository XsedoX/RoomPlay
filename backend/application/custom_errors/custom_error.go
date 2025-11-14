package custom_errors

import (
	"database/sql"
	"errors"
	"fmt"
)

type CustomError struct {
	Code      string
	ErrorType Type
	Err       error
	Title     string
}

func NewCustomError(code, title string, err error, errorType Type) *CustomError {
	if errors.Is(err, sql.ErrNoRows) {
		return &CustomError{Code: code, ErrorType: NotFound, Err: err, Title: title}
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
