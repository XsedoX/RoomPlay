package enqueued_song

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestNewEnqueuedSongSuccess(t *testing.T) {
	isrc := "USUM72005912"

	songData, err := song_data.NewSongData(
		faker.URL(),
		faker.Word(),
		faker.Name(),
		faker.URL(),
		5,
		music_provider.YouTube,
		&isrc,
	)

	enqueuedSong := NewEnqueuedSong(
		*songData,
		user_id.NewUserId(),
	)

	require.NoError(t, err)
	require.NotNil(t, enqueuedSong)
	require.WithinDuration(t, enqueuedSong.AddedAtUtc(), time.Now().UTC(), 1*time.Second)
	require.Nil(t, enqueuedSong.StartedAtUtc())
	require.Equal(t, int8(0), enqueuedSong.Votes())
	require.Equal(t, enqueued_song_state.Enqueued, enqueuedSong.State())
	require.NotNil(t, enqueuedSong.Id())
}
