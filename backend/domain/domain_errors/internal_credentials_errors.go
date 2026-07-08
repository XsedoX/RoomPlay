package domain_errors

func NewInternalCredentialsRefreshTokenEmptyError() error {
	return &DomainError{
		Code:        "InternalCredentials.RefreshToken.EmptyString",
		Description: "The field 'refresh token' cannot be an empty string.",
	}
}
