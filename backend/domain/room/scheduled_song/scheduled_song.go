package scheduled_song

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
)

type ScheduledSong struct {
	songData       song_data.SongData
	scheduledAtUtc time.Time
}

func (ss ScheduledSong) SongData() song_data.SongData {
	return ss.songData
}

func (ss ScheduledSong) ScheduledAtUtc() time.Time {
	return ss.scheduledAtUtc
}

func NewScheduledSong(songData song_data.SongData,
	scheduledAt time.Time,
) (*ScheduledSong, error) {
	if scheduledAt.Before(time.Now().UTC()) {
		return nil, domain_errors.NewScheduledSongScheduledInPastError()
	}
	return &ScheduledSong{
		songData:       songData,
		scheduledAtUtc: scheduledAt,
	}, nil
}
