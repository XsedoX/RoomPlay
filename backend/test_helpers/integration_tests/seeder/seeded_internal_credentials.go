package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
)

var (
	userSession = user_session.NewUserSession(users[0].Id(),
		users[0].Devices()[0].Id(),
	)
	internalCredentials = []internal_credentials.InternalCredentials{
		*internal_credentials.HydrateInternalCredentials(
			*userSession,
			"refreshTokenValue1",
			time.Now().AddDate(1, 0, 0),
			time.Date(2023, 12, 1, 12, 0o0, 0o0, 0o0, time.UTC)),
	}
)

func (s *Seeder) seedInternalCredentials(ctx context.Context, internalCredentials *internal_credentials.InternalCredentials) error {
	configuration := mock_configuration.MockConfiguration{}
	encrypter := encryper.NewEncrypter(configuration.Authentication().EncryptionKey)
	hashedRefreshToken := encrypter.Hash(internalCredentials.RefreshToken())
	_, err := s.Queryer.ExecContext(ctx, `
		INSERT INTO users_internal_credentials 
		(
			user_id,
			device_id,
			refresh_token,
			expires_at_utc,
			issued_at_utc
		)
		VALUES 
		(
			$1::uuid, $2::uuid, $3::bytea, $4, $5
		)
		`, internalCredentials.UserId().ToUuid(),
		internalCredentials.DeviceId().ToUuid(),
		hashedRefreshToken,
		internalCredentials.ExpiresAtUtc(),
		internalCredentials.IssuedAtUtc())
	if err != nil {
		return fmt.Errorf("failed to seed internal credentials: %w", err)
	}
	return nil
}
