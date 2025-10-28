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
func (r RefreshTokenRepository) UpdateToken(ctx context.Context, refreshToken *credentials.RefreshToken, queryer contracts.IQueryer) error {
	_, err := queryer.ExecContext(ctx,
		`
		UPDATE users_refresh_token SET 
			refresh_token=$1::bytea,
	    	expiration_timestamp=$2::timestamp,
		 	issued_at=$3::timestamp
		WHERE id = $4::uuid AND device_id = $5::uuid`,
		refreshToken.RefreshToken(),
		refreshToken.ExpirationTime().UTC(),
		refreshToken.IssuedAt().UTC(),
		uuid.UUID(refreshToken.Id()),
		uuid.UUID(refreshToken.DeviceId()))
	return err
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
	encryptedRefreshToken, err := r.encrypter.HashAndSalt(string(refreshToken.RefreshToken()))
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
			 expiration_timestamp,
			 issued_at
		)
		VALUES
		(
		     $1::uuid, $2::uuid, $3, $4, $5
		)`,
		uuid.UUID(refreshToken.Id()),
		uuid.UUID(refreshToken.DeviceId()),
		encryptedRefreshToken,
		refreshToken.ExpirationTime(),
		refreshToken.IssuedAt(),
	)
	return err
}
