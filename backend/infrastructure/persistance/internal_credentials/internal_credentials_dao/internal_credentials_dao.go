package internal_credentials_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/google/uuid"
)

type InternalCredentialsDao struct {
	UserId       uuid.UUID        `db:"user_id"`
	UserDao      user_dao.UserDao `gorm:"foreignKey:UserId;references:Id"`
	DeviceId     uuid.UUID        `db:"device_id"`
	RefreshToken []byte           `db:"refresh_token"`
	ExpiresAtUtc time.Time        `db:"expires_at_utc"`
	IssuedAtUtc  time.Time        `db:"issued_at_utc"`
}
