package domain_errors

func NewTokenExpiredError() error {
	return &DomainError{
		Code:        "Token.Expired",
		Description: "Token expiration time must be in the future",
	}
}

func NewTokenEmptyError() error {
	return &DomainError{
		Code:        "Token.Empty",
		Description: "Token cannot be empty",
	}
}
