package enqueued_song

import (
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewEnqueuedSongSuccess(t *testing.T) {
	songData, err := song_data.NewSongData(
		faker.URL(),
		faker.Word(),
		faker.Name(),
		faker.URL(),
		5,
	)
	roomId := room_id.NewRoomId()

	enqueuedSong := NewEnqueuedSong(
		*songData,
		user_id.NewUserId(),
		roomId,
	)

	assert.NoError(t, err)
	assert.NotNil(t, enqueuedSong)
	assert.WithinDuration(t, enqueuedSong.AddedAtUtc(), time.Now().UTC(), 1*time.Second)
	assert.Nil(t, enqueuedSong.StartedAtUtc())
	assert.Equal(t, uint8(0), enqueuedSong.Votes())
	assert.Equal(t, enqueued_song_state.Enqueued, enqueuedSong.State())
	assert.NotNil(t, enqueuedSong.Id())
}
