package internal_credentials

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/internal_credentials/user_session"
	"github.com/XsedoX/RoomPlay/domain/user/device/device_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewInternalCredentialsSuccess(t *testing.T) {
	userId := user_id.NewUserId()
	deviceId := device_id.NewDeviceId()
	internalCredentials := faker.UUIDDigit()
	userSession := user_session.NewUserSession(userId, deviceId)

	internalCredentialsObj, error := NewInternalCredentials(*userSession, internalCredentials)

	require.NoError(t, error)
	require.Equal(t, internalCredentials, internalCredentialsObj.RefreshToken())
	require.WithinDuration(t, internalCredentialsObj.IssuedAtUtc(), time.Now().UTC(), time.Second)
	require.WithinDuration(t, internalCredentialsObj.ExpiresAtUtc(), time.Now().Add(RefreshTokenExpirationTime).UTC(), time.Second)
}

func TestNewInternalCredentialsRefreshTokenEmpty(t *testing.T) {
	userId := user_id.NewUserId()
	deviceId := device_id.NewDeviceId()
	internalCredentials := ""
	userSession := user_session.NewUserSession(userId, deviceId)

	internalCredentialsObj, error := NewInternalCredentials(*userSession, internalCredentials)

	require.Error(t, error)
	require.Nil(t, internalCredentialsObj)
	castedError, ok := error.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "InternalCredentials.RefreshToken.EmptyString", castedError.Code)
	require.Equal(t, "The field 'refresh token' cannot be an empty string.", castedError.Description)
}
