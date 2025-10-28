package applicationErrors

import "errors"

type Type error

var (
	ErrValidation      = errors.New("validation error")
	ErrInfrastructure  = errors.New("infrastructure error")
	ErrAuthorization   = errors.New("authorization error")
	ErrUnexpectedError = errors.New("unexpected error")
)
