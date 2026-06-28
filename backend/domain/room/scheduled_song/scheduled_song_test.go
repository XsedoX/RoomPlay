package scheduled_song

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/time_before_now_domain_error"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestScheduledSongSuccess(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(300)
	scheduledAtUtc := time.Now().UTC().Add(10 * time.Minute)
	isrc := "USS1Z2500001"

	createdSongData, err := song_data.NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		music_provider.YouTube,
		&isrc,
	)
	require.NoError(t, err)

	scheduledSong, err := NewScheduledSong(
		*createdSongData,
		scheduledAtUtc,
	)

	require.NoError(t, err)
	require.NotNil(t, scheduledSong)
}

func TestScheduledSongScheduledBeforeNow(t *testing.T) {
	url := faker.URL()
	title := faker.Word()
	author := faker.Name()
	albumCoverUrl := faker.URL()
	lengthSeconds := uint16(300)
	scheduledAtUtc := time.Now().UTC().Add(-10 * time.Minute)
	musicProvider := music_provider.YouTube

	createdSongData, err := song_data.NewSongData(
		url,
		title,
		author,
		albumCoverUrl,
		lengthSeconds,
		musicProvider,
		nil,
	)
	require.NoError(t, err)

	scheduledSong, err := NewScheduledSong(
		*createdSongData,
		scheduledAtUtc,
	)

	require.Nil(t, scheduledSong)
	require.Error(t, err)
	castedErr, ok := err.(*time_before_now_domain_error.TimeBeforeNowDomainError)
	require.True(t, ok)
	require.Equal(t, "ScheduledSong.ScheduledAtUtc.TimeBeforeNow", castedErr.Code)
	require.Equal(t, "The field 'scheduled at' must be a time in the future.", castedErr.Description)
}
