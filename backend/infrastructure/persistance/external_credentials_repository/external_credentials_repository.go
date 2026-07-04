package external_credentials_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/dtos/refresh_access_token_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type ExternalCredentialsRepository struct {
	encrypter i_encrypter.IEncrypter
}

func NewExternalCredentialsRepository(encrypter i_encrypter.IEncrypter) *ExternalCredentialsRepository {
	return &ExternalCredentialsRepository{
		encrypter: encrypter,
	}
}

func (repo *ExternalCredentialsRepository) Grant(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error {
	encryptedAccessToken, err := repo.encrypter.Encrypt(credentials.AccessToken())
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := repo.encrypter.Encrypt(credentials.RefreshToken())
	if err != nil {
		return err
	}
	_, err = queryer.ExecContext(ctx,
		`
		INSERT INTO users_external_credentials 
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
		VALUES
		(
			$1::uuid, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (user_id) DO UPDATE
		SET external_id = EXCLUDED.external_id,
			  access_token = EXCLUDED.access_token,
			  refresh_token = EXCLUDED.refresh_token,
				music_provider = EXCLUDED.music_provider,
				access_token_expires_at_utc=EXCLUDED.access_token_expires_at_utc,
				refresh_token_expires_at_utc=EXCLUDED.refresh_token_expires_at_utc,
				issued_at_utc=EXCLUDED.issued_at_utc;
		`,
		credentials.Id().ToUuid(),
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

func (repo *ExternalCredentialsRepository) AccessTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (string, error) {
	encryptedAccessToken := make([]byte, 0)
	err := queryer.GetContext(ctx,
		&encryptedAccessToken,
		`
		select access_token::bytea
from users_external_credentials
where user_id = $1 and
access_token_expires_at_utc > (now()::timestamp + interval '1 minute');
  `, userId.ToUuid(),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	decryptedAccessToken, err := repo.encrypter.Decrypt(encryptedAccessToken)
	if err != nil {
		return "", err
	}
	return decryptedAccessToken, nil
}

func (repo *ExternalCredentialsRepository) RefreshAccessToken(ctx context.Context, refreshAccessTokenDto refresh_access_token_dto.RefreshAccessTokenDto, queryer i_queryer.IQueryer) error {
	encryptedAccessToken, err := repo.encrypter.Encrypt(refreshAccessTokenDto.AccessToken)
	if err != nil {
		return err
	}

	_, err = queryer.ExecContext(ctx,
		`
		UPDATE users_external_credentials
		SET access_token = $1, access_token_expires_at_utc = $2
		WHERE user_id = $3;
		`,
		encryptedAccessToken,
		refreshAccessTokenDto.AccessTokenExpiresAtUtc,
		refreshAccessTokenDto.UserId.ToUuid(),
	)
	if err != nil {
		return err
	}

	return nil
}
