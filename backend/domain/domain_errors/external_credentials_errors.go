package domain_errors

func NewExternalCredentialsRefreshTokenExpiredError() error {
	return &DomainError{
		Code:        "ExternalCredentials.RefreshToken.Expired",
		Description: "Refresh token expiration time must be in the future",
	}
}

func NewExternalCredentialsRefreshTokenEmptyError() error {
	return &DomainError{
		Code:        "ExternalCredentials.RefreshToken.Empty",
		Description: "Refresh token cannot be empty",
	}
}

func NewExternalCredentialsExternalIdEmptyError() error {
	return &DomainError{
		Code:        "ExternalCredentials.ExternalId.Empty",
		Description: "External id cannot be empty",
	}
}

func NewExternalCredentialsAccessTokenExpiredError() error {
	return &DomainError{
		Code:        "ExternalCredentials.AccessToken.Expired",
		Description: "Access token expiration time must be in the future",
	}
}

func NewExternalCredentialsAccessTokenEmptyError() error {
	return &DomainError{
		Code:        "ExternalCredentials.AccessToken.Empty",
		Description: "Access token cannot be empty",
	}
}
