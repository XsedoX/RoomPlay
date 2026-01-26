package persistance

import (
	"context"

	"github.com/google/uuid"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos"
)

type RefreshTokenRepository struct {
	encrypter application_contracts.IEncrypter
}

func NewRefreshTokenRepository(encrypter application_contracts.IEncrypter) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		encrypter: encrypter,
	}
}
func (r RefreshTokenRepository) GetTokenByValue(ctx context.Context, value string, queryer application_contracts.IQueryer) (*credentials.RefreshToken, error) {
	encryptedRefreshToken := r.encrypter.Hash(value)
	var tokenFromDb daos.RefreshTokenDao
	err := queryer.GetContext(ctx, &tokenFromDb,
		"SELECT * FROM users_refresh_tokens WHERE refresh_token = $1::bytea LIMIT 1;",
		encryptedRefreshToken)
	if err != nil {
		return nil, err
	}
	return credentials.HydrateRefreshToken(user.Id(tokenFromDb.UserId),
		user.DeviceId(tokenFromDb.DeviceId),
		string(tokenFromDb.RefreshToken),
		tokenFromDb.ExpiresAtUtc,
		tokenFromDb.IssuedAtUtc), nil
}
func (r RefreshTokenRepository) AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer application_contracts.IQueryer) error {
	encryptedRefreshToken := r.encrypter.Hash(refreshToken.RefreshToken())
	userId := refreshToken.Id()
	deviceId := refreshToken.DeviceId()
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
		refreshToken.ExpiresAtUtc(),
		refreshToken.IssuedAtUtc(),
	)
	return err
}
func (r RefreshTokenRepository) RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user.Id, deviceId *user.DeviceId, queryer application_contracts.IQueryer) error {
	uId := uuid.UUID(*userId)
	dId := uuid.UUID(*deviceId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2;",
		uId,
		dId)
	return err
}
func (r RefreshTokenRepository) RetireAllTokensByUserId(ctx context.Context, userId *user.Id, queryer application_contracts.IQueryer) error {
	id := uuid.UUID(*userId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_tokens WHERE user_id = $1;",
		id)
	return err
}
