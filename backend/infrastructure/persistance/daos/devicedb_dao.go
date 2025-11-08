package daos

import (
	"time"

	"github.com/google/uuid"
)

type DeviceDbDao struct {
	Id                uuid.UUID `db:"id"`
	FriendlyName      string    `db:"friendly_name"`
	IsHost            bool      `db:"is_host"`
	Type              string    `db:"type"`
	UserId            uuid.UUID `db:"user_id"`
	State             string    `db:"state"`
	LastLoggedInAtUtc time.Time `db:"last_logged_in_at_utc"`
}
