package persistance

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/google/uuid"
)

type ExternalCredentialsRepository struct {
	encrypter application_contracts.IEncrypter
}

func NewExternalCredentialsRepository(encrypter application_contracts.IEncrypter) *ExternalCredentialsRepository {
	return &ExternalCredentialsRepository{
		encrypter: encrypter,
	}
}

func (cr *ExternalCredentialsRepository) Grant(ctx context.Context, credentials *credentials.ExternalCredentials, queryer application_contracts.IQueryer) error {
	encryptedAccessToken, err := cr.encrypter.Encrypt(credentials.AccessToken())
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := cr.encrypter.Encrypt(credentials.RefreshToken())
	if err != nil {
		return err
	}
	_, err = queryer.ExecContext(ctx,
		`INSERT INTO users_external_credentials 
    			(
						user_id, 
		external_id,
						access_token, 
						refresh_token, 
		music_provider,
						access_token_expires_at_utc, 
						refresh_token_expires_at_utc,
    			 		issued_at_utc
				) 
				VALUES (
				        $1::uuid, $2, $3, $4, $5, $6, $7, $8
				)
				ON CONFLICT (user_id) DO UPDATE
				SET external_id = EXCLUDED.external_id,
		access_token = EXCLUDED.access_token,
				    refresh_token = EXCLUDED.refresh_token,
		music_provider = EXCLUDED.music_provider,
				    access_token_expires_at_utc=EXCLUDED.access_token_expires_at_utc,
				    refresh_token_expires_at_utc=EXCLUDED.refresh_token_expires_at_utc,
				    issued_at_utc=EXCLUDED.issued_at_utc;`,
		uuid.UUID(credentials.Id()),
		credentials.ExternalId(),
		encryptedAccessToken,
		encryptedRefreshToken,
		credentials.MusicProvider().String(),
		credentials.AccessTokenExpiresAtUtc(),
		credentials.RefreshTokenExpiresAtUtc(),
		credentials.IssuedAtUtc(),
	)
	return err
}

func (credentials *ExternalCredentialsRepository) GetAccessTokenByUserId(ctx context.Context, userId user.Id, queryer application_contracts.IQueryer) (string, error) {
	return "", nil
}
