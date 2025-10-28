package persistance

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/domain/credentials"
)

type CredentialsRepository struct {
	encrypter contracts.IEncrypter
}

func NewCredentialsRepository(encrypter contracts.IEncrypter) *CredentialsRepository {
	return &CredentialsRepository{
		encrypter: encrypter,
	}
}
func (cr *CredentialsRepository) Grant(ctx context.Context, credentials *credentials.External, queryer contracts.IQueryer) error {
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
						access_token_expiration_timestamp, 
						refresh_token_expiration_timestamp,
    			 		issued_at
				) 
				VALUES (
				        $1::uuid, $2, $3, $4, $5, $6, $7
				)`,
		uuid.UUID(credentials.Id()),
		encryptedAccessToken,
		encryptedRefreshToken,
		strings.Join(credentials.Scopes(), " "),
		credentials.AccessTokenExpirationTime(),
		credentials.RefreshTokenExpirationTime(),
		credentials.IssuedAt(),
	)
	return err
}
