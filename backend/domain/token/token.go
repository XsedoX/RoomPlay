package token

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
)

type Token struct {
	value        string
	expiresAtUtc time.Time
}

func NewToken(value string, expiresAtUtc time.Time) (*Token, error) {
	if value == "" {
		return nil, domain_errors.NewTokenEmptyError()
	}
	if expiresAtUtc.Before(time.Now().UTC()) {
		return nil, domain_errors.NewTokenExpiredError()
	}

	return &Token{
		value:        value,
		expiresAtUtc: expiresAtUtc.UTC(),
	}, nil
}

func (t Token) Value() string {
	return t.value
}

func (t Token) ExpiresAtUtc() time.Time {
	return t.expiresAtUtc.UTC()
}

func HydrateToken(value string, expiresAtUtc time.Time) *Token {
	return &Token{
		value:        value,
		expiresAtUtc: expiresAtUtc.UTC(),
	}
}

func (t Token) IsExpired() bool {
	return t.expiresAtUtc.Before(time.Now().UTC())
}
