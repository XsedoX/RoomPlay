package value_objects

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func NewEmail(value string) (*Email, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return nil, errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(trimmedValue) {
		return nil, errors.New("invalid email format")
	}
	return &Email{value: trimmedValue}, nil
}
func (e Email) Value() string {
	return e.value
}
func (e Email) Equal(other Email) bool {
	return e.value == other.value
}
