package persistance

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
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
	var tokenFromDb []credentials.RefreshToken
	err := queryer.SelectContext(ctx, &tokenFromDb,
		"SELECT * FROM users_refresh_token WHERE refresh_token = ?::bytea;",
		[]byte(value))
	if err != nil {
		return nil, err
	}
	if len(tokenFromDb) == 0 {
		return nil, errors.New("refresh token not found")
	}
	return &tokenFromDb[0], nil
}
func (r RefreshTokenRepository) AssignNewToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer contracts.IQueryer) error {
	encryptedRefreshToken, err := r.encrypter.HashAndSalt(refreshToken.RefreshToken())
	if err != nil {
		return err
	}
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
		uuid.UUID(refreshToken.Id()),
		uuid.UUID(refreshToken.DeviceId()),
		encryptedRefreshToken,
		refreshToken.ExpiresAtUtc(),
		refreshToken.IssuedAtUtc(),
	)
	return err
}
