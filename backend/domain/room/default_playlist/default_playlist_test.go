package default_playlist

import (
	"testing"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
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
	castErr, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.ExternalId.EmptyString", castErr.Code)
	require.Equal(t, "The field 'external id' cannot be an empty.", castErr.Description)

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
	titleError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.Title.EmptyString", titleError.Code)
	require.Equal(t, "The field 'title' cannot be an empty.", titleError.Description)
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
	castErr, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "DefaultPlaylist.SongsAmount.Zero", castErr.Code)
	require.Equal(t, "Playlist has to have at least one song.", castErr.Description)
}
