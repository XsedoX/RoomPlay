package shared

import "github.com/google/uuid"

type UserId uuid.UUID

func (id UserId) String() string {
	return uuid.UUID(id).String()
}
