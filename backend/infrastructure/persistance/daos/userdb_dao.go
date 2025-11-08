package daos

import (
	"time"

	"github.com/google/uuid"
)

type UserDbDao struct {
	Id             uuid.UUID  `db:"id"`
	ExternalId     string     `db:"external_id"`
	Name           string     `db:"name"`
	Surname        string     `db:"surname"`
	RoomId         *uuid.UUID `db:"room_id"`
	Role           *string    `db:"role"`
	BoostUsedAtUtc *time.Time `db:"used_at_utc"`
}
