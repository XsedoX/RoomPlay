package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
	"github.com/go-faker/faker/v4"
)

var (
	externalCredentialsExpiration1 = time.Now().Add(time.Hour * 1).UTC()
	externalCredentialsExpiration2 = time.Now().Add(time.Hour * 2).UTC()
	externalCredentialsExpiration3 = time.Now().Add(time.Hour * 2).UTC()
	externalCredentialsExpiration4 = time.Now().Add(time.Hour * 3).UTC()
	externalCredentialsExpiration5 = time.Now().Add(time.Hour * 4).UTC()
	refreshCredentialsExpiration1  = time.Now().AddDate(0, 1, 0).UTC()
	refreshCredentialsExpiration2  = time.Now().AddDate(0, 2, 0).UTC()
	refreshCredentialsExpiration3  = time.Now().AddDate(0, 2, 0).UTC()
	refreshCredentialsExpiration4  = time.Now().AddDate(0, 3, 0).UTC()
	refreshCredentialsExpiration5  = time.Now().AddDate(0, 4, 0).UTC()
	externalCredentials            = []external_credentials.ExternalCredentials{
		*external_credentials.HydrateExternalCredentials(users[0].Id(),
			faker.Jwt(),
			faker.Jwt(),
			faker.UUIDDigit(),
			music_provider.YouTube,
			externalCredentialsExpiration1,
			refreshCredentialsExpiration1,
			time.Now().UTC(),
		),
		*external_credentials.HydrateExternalCredentials(users[1].Id(),
			faker.Jwt(),
			faker.Jwt(),
			faker.UUIDDigit(),
			music_provider.Spotify,
			externalCredentialsExpiration2,
			refreshCredentialsExpiration2,
			time.Now().UTC(),
		),
		*external_credentials.HydrateExternalCredentials(users[2].Id(),
			faker.Jwt(),
			faker.Jwt(),
			faker.UUIDDigit(),
			music_provider.Spotify,
			externalCredentialsExpiration3,
			refreshCredentialsExpiration3,
			time.Now().UTC(),
		),
		*external_credentials.HydrateExternalCredentials(users[3].Id(),
			faker.Jwt(),
			faker.Jwt(),
			faker.UUIDDigit(),
			music_provider.YouTube,
			externalCredentialsExpiration4,
			refreshCredentialsExpiration4,
			time.Now().UTC(),
		),
		*external_credentials.HydrateExternalCredentials(users[4].Id(),
			faker.Jwt(),
			faker.Jwt(),
			faker.UUIDDigit(),
			music_provider.Spotify,
			externalCredentialsExpiration5,
			refreshCredentialsExpiration5,
			time.Now().UTC(),
		),
	}
)

func (s *Seeder) seedExternalCredentials(ctx context.Context, creds *external_credentials.ExternalCredentials) error {
	configuration := mock_configuration.MockConfiguration{}
	encrypter := encryper.NewEncrypter(&configuration)

	encryptedAccessToken, _ := encrypter.Encrypt(creds.AccessToken())
	encryptedRefreshToken, _ := encrypter.Encrypt(creds.RefreshToken())
	_, err := s.Queryer.ExecContext(ctx, `
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
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		`,
		creds.Id(),
		creds.ExternalId(),
		encryptedAccessToken,
		encryptedRefreshToken,
		creds.MusicProvider().String(),
		creds.AccessTokenExpiresAtUtc(),
		creds.RefreshTokenExpiresAtUtc(),
		creds.IssuedAtUtc(),
	)
	if err != nil {
		return fmt.Errorf("failed to seed external credentials: %w", err)
	}
	return nil
}
