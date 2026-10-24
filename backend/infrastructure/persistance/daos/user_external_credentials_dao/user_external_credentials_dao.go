package user_external_credentials_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/google/uuid"
)

type UserExternalCredentialsDao struct {
	UserId                  uuid.UUID        `gorm:"primaryKey;foreignKey:UserId;references:Id"`
	UserDao                 user_dao.UserDao `gorm:"foreignKey:UserId;references:Id"`
	ExternalId              string           `gorm:"column:external_id"`
	AccessToken             []byte           `gorm:"column:access_token;type:bytea"`
	RefreshToken            []byte           `gorm:"column:refresh_token;type:bytea"`
	MusicProvider           string           `gorm:"column:music_provider;type:music_provider"`
	AccessTokenExpiresAtUtc time.Time        `gorm:"column:access_token_expires_at_utc;type:timestamp"`
	AccessTokenIssuedAtUtc  time.Time        `gorm:"column:access_token_issued_at_utc;type:timestamp"`
}

func (UserExternalCredentialsDao) TableName() string {
	return "users_external_credentials"
}

// ToDomain decrypts the stored tokens; hydration also supports expired tokens.
func (d UserExternalCredentialsDao) ToDomain(encrypter i_encrypter.IEncrypter) (*external_credentials.ExternalCredentials, error) {
	accessToken, err := encrypter.Decrypt(d.AccessToken)
	if err != nil {
		return nil, err
	}
	refreshToken, err := encrypter.Decrypt(d.RefreshToken)
	if err != nil {
		return nil, err
	}
	return external_credentials.HydrateExternalCredentials(
		user_id.UserId(d.UserId),
		accessToken,
		refreshToken,
		d.ExternalId,
		*music_provider.ParseMusicProvider(d.MusicProvider),
		d.AccessTokenExpiresAtUtc,
		d.AccessTokenIssuedAtUtc,
	), nil
}

func (d *UserExternalCredentialsDao) FromDomain(value *external_credentials.ExternalCredentials, encrypter i_encrypter.IEncrypter) error {
	accessToken, err := encrypter.Encrypt(value.GetAccessToken().Value())
	if err != nil {
		return err
	}
	refreshToken, err := encrypter.Encrypt(value.GetRefreshToken())
	if err != nil {
		return err
	}
	*d = UserExternalCredentialsDao{
		UserId:                  value.Id().ToUuid(),
		ExternalId:              value.ExternalId(),
		AccessToken:             accessToken,
		RefreshToken:            refreshToken,
		MusicProvider:           value.MusicProvider().String(),
		AccessTokenExpiresAtUtc: value.AccessTokenExpiresAtUtc(),
		AccessTokenIssuedAtUtc:  value.IssuedAtUtc(),
	}
	return nil
}
