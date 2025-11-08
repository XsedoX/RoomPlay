package persistance

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/infrastructure/persistance/daos"
)

type RefreshTokenRepository struct {
	encrypter contracts.IEncrypter
}

func NewRefreshTokenRepository(encrypter contracts.IEncrypter) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		encrypter: encrypter,
	}
}
func (r RefreshTokenRepository) GetTokenByValue(ctx context.Context, value string, queryer contracts.IQueryer) (*credentials.RefreshToken, error) {
	var tokenFromDb []daos.RefreshTokenDbDao
	err := queryer.SelectContext(ctx, &tokenFromDb,
		"SELECT * FROM users_refresh_token WHERE refresh_token = $1::bytea;",
		[]byte(value))
	if err != nil {
		return nil, err
	}
	if len(tokenFromDb) == 0 {
		return nil, errors.New("refresh token not found")
	}
	return credentials.HydrateRefreshToken(user.Id(tokenFromDb[0].UserId),
		user.DeviceId(tokenFromDb[0].DeviceId),
		string(tokenFromDb[0].RefreshToken),
		tokenFromDb[0].ExpiresAtUtc,
		tokenFromDb[0].IssuedAtUtc), nil
}
func (r RefreshTokenRepository) AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer contracts.IQueryer) error {
	encryptedRefreshToken, err := r.encrypter.HashAndSalt(refreshToken.RefreshToken())
	if err != nil {
		return err
	}
	userId := refreshToken.Id()
	deviceId := refreshToken.DeviceId()
	_, err = queryer.ExecContext(ctx,
		`
		INSERT INTO users_refresh_token 
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

func (r RefreshTokenRepository) RetireTokenByUserIdAndDeviceId(ctx context.Context, userId *user.Id, deviceId *user.DeviceId, queryer contracts.IQueryer) error {
	uId := uuid.UUID(*userId)
	dId := uuid.UUID(*deviceId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_token WHERE user_id = $1 AND device_id = $2;",
		uId,
		dId)
	return err
}
func (r RefreshTokenRepository) RetireTokenByUserId(ctx context.Context, userId *user.Id, queryer contracts.IQueryer) error {
	id := uuid.UUID(*userId)
	_, err := queryer.ExecContext(ctx,
		"DELETE FROM users_refresh_token WHERE user_id = $1;",
		id)
	return err
}
