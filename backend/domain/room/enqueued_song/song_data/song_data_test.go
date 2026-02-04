package song_data

import (
	"testing"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewSongDataSuccess(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(300)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.NoError(t, err)
	require.NotNil(t, createdSongData)
}

func TestNewSongDataLengthZero(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(0)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*validation_domain_error.ValidationDomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.LengthSeconds.Zero", castedError.Code)
	require.Equal(t, "Song length in seconds cannot be zero", castedError.Description)
}

func TestNewSongDataUrlEmpty(t *testing.T) {
	url := ""
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.URL.EmptyString", castedError.Code)
	require.Equal(t, "The field 'url' cannot be an empty string.", castedError.Description)
}

func TestNewSongDataTitleEmpty(t *testing.T) {
	url := faker.URL()
	title := ""
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.Title.EmptyString", castedError.Code)
	require.Equal(t, "The field 'title' cannot be an empty string.", castedError.Description)
}

func TestNewSongDataAuthorEmpty(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := ""
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(5)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.Author.EmptyString", castedError.Code)
	require.Equal(t, "The field 'author' cannot be an empty string.", castedError.Description)
}

func TestNewSongDataAlbumCoverUrlEmpty(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := ""
	lengthSeconds := uint16(5)

	createdSongData, err := NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
	)

	require.Error(t, err)
	require.Nil(t, createdSongData)
	castedError, ok := err.(*empty_string_domain_error.EmptyStringDomainError)
	require.True(t, ok)
	require.Equal(t, "SongData.AlbumCoverURL.EmptyString", castedError.Code)
	require.Equal(t, "The field 'album cover url' cannot be an empty string.", castedError.Description)
}
