package internal_credentials_dao

import (
	"time"

	"github.com/google/uuid"
)

type InternalCredentialsDao struct {
	UserId       uuid.UUID `db:"user_id"`
	DeviceId     uuid.UUID `db:"device_id"`
	RefreshToken []byte    `db:"refresh_token"`
	ExpiresAtUtc time.Time `db:"expires_at_utc"`
	IssuedAtUtc  time.Time `db:"issued_at_utc"`
}
