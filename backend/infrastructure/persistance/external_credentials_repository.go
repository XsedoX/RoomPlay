package persistance

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
)

type ExternalCredentialsRepository struct {
	encrypter contracts.IEncrypter
}

func NewExternalCredentialsRepository(encrypter contracts.IEncrypter) *ExternalCredentialsRepository {
	return &ExternalCredentialsRepository{
		encrypter: encrypter,
	}
}
func (cr *ExternalCredentialsRepository) Grant(ctx context.Context, credentials *credentials.External, queryer contracts.IQueryer) error {
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
						access_token, 
						refresh_token, 
						scope, 
						access_token_expires_at_utc, 
						refresh_token_expires_at_utc,
    			 		issued_at_utc
				) 
				VALUES (
				        $1::uuid, $2, $3, $4, $5, $6, $7
				)
				ON CONFLICT (user_id) DO UPDATE
				SET access_token = EXCLUDED.access_token,
				    refresh_token = EXCLUDED.refresh_token,
				    scope = EXCLUDED.scope,
				    access_token_expires_at_utc=EXCLUDED.access_token_expires_at_utc,
				    refresh_token_expires_at_utc=EXCLUDED.refresh_token_expires_at_utc,
				    issued_at_utc=EXCLUDED.issued_at_utc;`,
		uuid.UUID(credentials.Id()),
		encryptedAccessToken,
		encryptedRefreshToken,
		strings.Join(credentials.Scopes(), " "),
		credentials.AccessTokenExpiresAtUtc(),
		credentials.RefreshTokenExpiresAtUtc(),
		credentials.IssuedAtUtc(),
	)
	return err
}
