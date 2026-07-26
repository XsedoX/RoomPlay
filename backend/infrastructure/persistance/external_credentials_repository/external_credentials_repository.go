package external_credentials_repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_queryer"
	"github.com/XsedoX/RoomPlay/application/dtos/refresh_access_token_dto"
	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/token"
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
	encryptedAccessToken, err := repo.encrypter.Encrypt(credentials.GetAccessToken().Value())
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := repo.encrypter.Encrypt(credentials.GetRefreshToken().Value())
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

func (repo *ExternalCredentialsRepository) GetExternalCredentialsByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*external_credentials.ExternalCredentials, error) {
	type result struct {
		AccessToken              []byte    `db:"access_token"`
		RefreshToken             []byte    `db:"refresh_token"`
		ExternalId               string    `db:"external_id"`
		MusicProvider            string    `db:"music_provider"`
		AccessTokenExpiresAtUtc  time.Time `db:"access_token_expires_at_utc"`
		RefreshTokenExpiresAtUtc time.Time `db:"refresh_token_expires_at_utc"`
		IssuedAtUtc              time.Time `db:"issued_at_utc"`
	}
	var resultInstance result
	err := queryer.GetContext(ctx,
		&resultInstance,
		`
		select access_token::bytea,
			refresh_token::bytea,
			external_id,
			music_provider,
			access_token_expires_at_utc,
			refresh_token_expires_at_utc,
			issued_at_utc
		from users_external_credentials 
		where user_id = $1;
		`,
		userId.ToUuid(),
	)
	if err != nil {
		return nil, err
	}
	decryptedAccessToken, err := repo.encrypter.Decrypt(resultInstance.AccessToken)
	if err != nil {
		return nil, err
	}
	decryptedRefreshToken, err := repo.encrypter.Decrypt(resultInstance.RefreshToken)
	if err != nil {
		return nil, err
	}
	return external_credentials.HydrateExternalCredentials(
		userId,
		decryptedAccessToken,
		decryptedRefreshToken,
		resultInstance.ExternalId,
		*music_provider.ParseMusicProvider(resultInstance.MusicProvider),
		resultInstance.AccessTokenExpiresAtUtc,
		resultInstance.RefreshTokenExpiresAtUtc,
		resultInstance.IssuedAtUtc,
	), nil
}

func (repo *ExternalCredentialsRepository) UpdateExternalCredentials(ctx context.Context, credentials *external_credentials.ExternalCredentials, queryer i_queryer.IQueryer) error {
	encryptedAccessToken, err := repo.encrypter.Encrypt(credentials.GetAccessToken().Value())
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := repo.encrypter.Encrypt(credentials.GetRefreshToken().Value())
	if err != nil {
		return err
	}
	_, err = queryer.ExecContext(ctx,
		`
		UPDATE users_external_credentials
		SET access_token = $1,
			refresh_token = $2,
			external_id = $3,
			music_provider = $4,
			access_token_expires_at_utc = $5,
			refresh_token_expires_at_utc = $6,
			issued_at_utc = $7
		WHERE user_id = $8;
		`,
		encryptedAccessToken,
		encryptedRefreshToken,
		credentials.ExternalId(),
		credentials.MusicProvider().String(),
		credentials.AccessTokenExpiresAtUtc(),
		credentials.RefreshTokenExpiresAtUtc(),
		credentials.IssuedAtUtc(),
		credentials.Id().ToUuid(),
	)
	return err
}

func (repo *ExternalCredentialsRepository) GetRefreshTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error) {
	type result struct {
		RefreshToken             []byte    `db:"refresh_token"`
		RefreshTokenExpiresAtUtc time.Time `db:"refresh_token_expires_at_utc"`
	}
	var resultInstance result
	err := queryer.GetContext(ctx,
		&resultInstance,
		`
		select refresh_token::bytea, refresh_token_expires_at_utc
from users_external_credentials
where user_id = $1;
  `,
		userId.ToUuid(),
		&resultInstance,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}
	decryptedRefreshToken, err := repo.encrypter.Decrypt(resultInstance.RefreshToken)
	if err != nil {
		return nil, err
	}
	return token.HydrateToken(decryptedRefreshToken, resultInstance.RefreshTokenExpiresAtUtc), nil
}

func (repo *ExternalCredentialsRepository) GetAccessTokenByUserId(ctx context.Context, userId user_id.UserId, queryer i_queryer.IQueryer) (*token.Token, error) {
	type result struct {
		AccessToken             []byte    `db:"access_token"`
		AccessTokenExpiresAtUtc time.Time `db:"access_token_expires_at_utc"`
	}
	var resultInstance result
	err := queryer.GetContext(ctx,
		&resultInstance,
		`
		select access_token::bytea, access_token_expires_at_utc
from users_external_credentials
where user_id = $1;
  `, userId.ToUuid(),
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	decryptedAccessToken, err := repo.encrypter.Decrypt(resultInstance.AccessToken)
	if err != nil {
		return nil, err
	}
	return token.HydrateToken(decryptedAccessToken, resultInstance.AccessTokenExpiresAtUtc), nil
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
