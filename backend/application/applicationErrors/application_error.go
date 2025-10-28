package applicationErrors

type ApplicationError struct {
	Message   string
	ErrorType Type
	Err       error
}

func (e ApplicationError) Error() string {
	return e.Message
}
func NewApplicationError(message string, err error, errType Type) *ApplicationError {
	return &ApplicationError{Message: message, Err: err, ErrorType: errType}
}
