package credentials

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshToken(t *testing.T) {
	userId := user.Id(uuid.New())
	deviceId := user.DeviceId(uuid.New())
	refreshToken := faker.UUIDDigit()

	refreshTokenObj := NewRefreshToken(userId, deviceId, refreshToken)

	require.Equal(t, userId, refreshTokenObj.Id())
	require.Equal(t, deviceId, refreshTokenObj.DeviceId())
	require.Equal(t, refreshToken, refreshTokenObj.RefreshToken())
	require.WithinDuration(t, refreshTokenObj.IssuedAtUtc(), time.Now().UTC(), time.Second)
	require.WithinDuration(t, refreshTokenObj.ExpiresAtUtc(), time.Now().Add(RefreshTokenExpirationTime).UTC(), time.Second)
}
