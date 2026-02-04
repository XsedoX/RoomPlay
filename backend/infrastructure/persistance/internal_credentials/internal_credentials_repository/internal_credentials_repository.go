package internal_credentials_repository

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_dao"
	"github.com/google/uuid"
)

type InternalCredentialsRepository struct {
	encrypter i_encrypter.IEncrypter
}

func NewInternalCredentialsRepository(encrypter i_encrypter.IEncrypter) *InternalCredentialsRepository {
	return &InternalCredentialsRepository{
		encrypter: encrypter,
	}
}

func (r InternalCredentialsRepository) GetTokenByValue(ctx context.Context, value string, queryer i_queryer.IQueryer) (*internal_credentials.InternalCredentials, error) {
	encryptedRefreshToken := r.encrypter.Hash(value)
	var tokenFromDb internal_credentials_dao.InternalCredentialsDao
	err := queryer.GetContext(ctx, &tokenFromDb,
		"SELECT * FROM users_refresh_tokens WHERE refresh_token = $1::bytea LIMIT 1;",
		encryptedRefreshToken)
	if err != nil {
		return nil, err
	}
	userSession := user_session.NewUserSession(user_id.UserId(tokenFromDb.UserId), device_id.DeviceId(tokenFromDb.DeviceId))
	return internal_credentials.HydrateInternalCredentials(
		*userSession,
		string(tokenFromDb.RefreshToken),
		tokenFromDb.ExpiresAtUtc,
		tokenFromDb.IssuedAtUtc), nil
}

func (r InternalCredentialsRepository) AssignNewToken(ctx context.Context, internaCredentials *internal_credentials.InternalCredentials, queryer i_queryer.IQueryer) error {
	encryptedRefreshToken := r.encrypter.Hash(internaCredentials.RefreshToken())
	userSession := internaCredentials.Id()
	userId := userSession.UserId()
	deviceId := userSession.DeviceId()
	_, err := queryer.ExecContext(ctx,
		`
		INSERT INTO users_refresh_tokens 
		(
			 user_id,
		 	 device_id,
			 refresh_token,
			 expires_at_utc,
			 issued_at_utc
		)
		VALUES
		(
		     $1::uuid, $2::uuid, $3, $4, $5
		)
		ON CONFLICT (user_id, device_id) DO UPDATE
		SET 
		    refresh_token = EXCLUDED.refresh_token,
		    expires_at_utc = EXCLUDED.expires_at_utc,
		    issued_at_utc = EXCLUDED.issued_at_utc;`,
		userId.ToUuid(),
		deviceId.ToUuid(),
		encryptedRefreshToken,
		internaCredentials.ExpiresAtUtc(),
		internaCredentials.IssuedAtUtc(),
	)
	return err
}

func (r InternalCredentialsRepository) RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user_id.UserId, deviceId *device_id.DeviceId, queryer i_queryer.IQueryer) error {
	uId := uuid.UUID(*userId)
	dId := uuid.UUID(*deviceId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2;",
		uId,
		dId)
	return err
}

func (r InternalCredentialsRepository) RetireAllTokensByUserId(ctx context.Context, userId *user_id.UserId, queryer i_queryer.IQueryer) error {
	id := uuid.UUID(*userId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_tokens WHERE user_id = $1;",
		id)
	return err
}
