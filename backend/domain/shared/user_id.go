package shared

import "github.com/google/uuid"

type UserId uuid.UUID

func (id UserId) String() string {
	return uuid.UUID(id).String()
}

// ParseUserId returns nil if something went wrong
func ParseUserId(value string) *UserId {
	result, err := uuid.Parse(value)
	if err != nil {
		return nil
	}
	userIdRes := UserId(result)
	return &userIdRes
}
