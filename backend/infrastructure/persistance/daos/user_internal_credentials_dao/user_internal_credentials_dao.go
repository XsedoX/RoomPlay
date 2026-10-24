package user_internal_credentials_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/device_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/google/uuid"
)

type UserInternalCredentialsDao struct {
	UserId       uuid.UUID            `gorm:"primaryKey;foreignKey:UserId;references:Id"`
	UserDao      user_dao.UserDao     `gorm:"foreignKey:UserId;references:Id"`
	DeviceId     uuid.UUID            `gorm:"primaryKey;foreignKey:DeviceId;references:Id"`
	DeviceDao    device_dao.DeviceDao `gorm:"foreignKey:DeviceId;references:Id"`
	RefreshToken []byte               `gorm:"column:refresh_token;type:bytea"`
	ExpiresAtUtc time.Time            `gorm:"column:expires_at_utc;type:timestamp"`
	IssuedAtUtc  time.Time            `gorm:"column:issued_at_utc;type:timestamp"`
}

func (UserInternalCredentialsDao) TableName() string {
	return "users_internal_credentials"
}

// ToDomain requires the plaintext token whose hash was used to find this row.
// RefreshToken contains a one-way hash, so the DAO cannot recover its plaintext.
// Token verification remains the caller's responsibility.
func (d UserInternalCredentialsDao) ToDomain(refreshToken string) *internal_credentials.InternalCredentials {
	session := user_session.NewUserSession(user_id.UserId(d.UserId), device_id.DeviceId(d.DeviceId))
	return internal_credentials.HydrateInternalCredentials(
		*session,
		refreshToken,
		d.ExpiresAtUtc,
		d.IssuedAtUtc,
	)
}

func (d *UserInternalCredentialsDao) FromDomain(value *internal_credentials.InternalCredentials, encrypter i_encrypter.IEncrypter) {
	*d = UserInternalCredentialsDao{
		UserId:       value.UserId().ToUuid(),
		DeviceId:     value.DeviceId().ToUuid(),
		RefreshToken: encrypter.Hash(value.RefreshToken().Value()),
		ExpiresAtUtc: value.ExpiresAtUtc(),
		IssuedAtUtc:  value.IssuedAtUtc(),
	}
}
