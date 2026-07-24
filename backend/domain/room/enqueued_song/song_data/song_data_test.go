package song_data

import (
	"testing"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewSongDataSuccess(t *testing.T) {
	externalId := faker.UUIDDigit()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(300)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.NoError(t, err)
	require.NotNil(t, createdSongData)
}

func TestNewSongDataLengthZero(t *testing.T) {
	externalId := faker.UUIDDigit()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(0)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.LengthSeconds.Zero", castedError.Code)
	require.Equal(t, "Song length in seconds cannot be zero", castedError.Description)
}

func TestNewSongDataExternalIdEmpty(t *testing.T) {
	externalId := ""
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.ExternalId.Empty", castedError.Code)
	require.Equal(t, "Song external id cannot be empty", castedError.Description)
}

func TestNewSongDataTitleEmpty(t *testing.T) {
	externalId := faker.UUIDDigit()
	title := ""
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.Title.Empty", castedError.Code)
	require.Equal(t, "Song title cannot be empty", castedError.Description)
}

func TestNewSongDataAuthorEmpty(t *testing.T) {
	externalId := faker.UUIDDigit()
	title := faker.Word()
	author := ""
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.Author.Empty", castedError.Code)
	require.Equal(t, "Song author cannot be empty", castedError.Description)
}

func TestNewSongDataAlbumCoverUrlEmpty(t *testing.T) {
	externalId := faker.UUIDDigit()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := ""
	lengthSeconds := uint16(5)
	musicProvider := music_provider.YouTube
	isrc := "USS1Z2500001"

	createdSongData, err := NewSongData(
		externalId,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		&isrc,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*domain_errors.DomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.AlbumCoverUrl.Empty", castedError.Code)
	require.Equal(t, "Song album cover url cannot be empty", castedError.Description)
}
