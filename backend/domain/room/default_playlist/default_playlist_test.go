package default_playlist

import (
	"testing"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultPlaylistSuccess(t *testing.T) {
	externalId := faker.UUIDDigit()
	userId := user_id.NewUserId()
	songsAmount := uint(10)
	title := faker.Word()

	playlist, err := NewDefaultPlaylist(
		externalId,
		userId,
		songsAmount,
		title,
	)

	require.NoError(t, err)
	require.NotNil(t, playlist)
	require.Equal(t, externalId, playlist.ExternalId())
	require.Equal(t, userId, playlist.UserId())
	require.Equal(t, songsAmount, playlist.SongsAmount())
	require.Equal(t, title, playlist.Title())
}

func TestNewDefaultPlaylistEmptyString(t *testing.T) {
	userId := user_id.NewUserId()
	songsAmount := uint(10)
	title := faker.Word()

	// externalId empty string
	externalId := ""

	playlist, err := NewDefaultPlaylist(
		externalId,
		userId,
		songsAmount,
		title,
	)

	require.Error(t, err)
	require.Nil(t, playlist)
	castErr, ok := err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.ExternalId.EmptyString", castErr.Code)
	require.Equal(t, "The field 'external id' cannot be an empty string.", castErr.Description)

	// titile empty string
	externalId = faker.UUIDDigit()
	title = ""
	playlist, err = NewDefaultPlaylist(
		externalId,
		userId,
		songsAmount,
		title,
	)

	require.Error(t, err)
	require.Nil(t, playlist)
	castErr, ok = err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.Title.EmptyString", castErr.Code)
	require.Equal(t, "The field 'title' cannot be an empty string.", castErr.Description)
}

func TestNewDefaultPlaylistSongsAmountZero(t *testing.T) {
	userId := user_id.NewUserId()
	songsAmount := uint(0)
	title := faker.Word()
	externalId := faker.UUIDDigit()

	playlist, err := NewDefaultPlaylist(
		externalId,
		userId,
		songsAmount,
		title,
	)

	require.Error(t, err)
	require.Nil(t, playlist)
	castErr, ok := err.(*validation_domain_error.ValidationDomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.SongsAmount.Zero", castErr.Code)
	require.Equal(t, "Songs amount cannot be zero", castErr.Description)
}
