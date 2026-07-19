package enqueued_song

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type EnqueuedSong struct {
	shared.Entity[enqueued_song_id.EnqueuedSongId]
	songData     song_data.SongData
	addedAtUtc   time.Time
	startedAtUtc *time.Time
	state        enqueued_song_state.EnqueuedSongState
	votes        int8
	addedBy      user_id.UserId
}

func (s EnqueuedSong) AddedBy() user_id.UserId {
	return s.addedBy
}

func (s EnqueuedSong) AddedAtUtc() time.Time {
	return s.addedAtUtc
}

func (s EnqueuedSong) StartedAtUtc() *time.Time {
	return s.startedAtUtc
}

func (s EnqueuedSong) State() enqueued_song_state.EnqueuedSongState {
	return s.state
}

func (s EnqueuedSong) Votes() int8 {
	return s.votes
}

func (s EnqueuedSong) SongData() song_data.SongData {
	return s.songData
}

func (s *EnqueuedSong) IsPlaying() bool {
	return s.state == enqueued_song_state.Playing
}

func HydrateEnqueuedSong(
	id enqueued_song_id.EnqueuedSongId,
	songData song_data.SongData,
	addedAtUtc time.Time,
	startedAtUtc *time.Time,
	state enqueued_song_state.EnqueuedSongState,
	votes int8,
	addedBy user_id.UserId,
) *EnqueuedSong {
	result := &EnqueuedSong{
		songData:     songData,
		addedAtUtc:   addedAtUtc,
		startedAtUtc: startedAtUtc,
		state:        state,
		votes:        votes,
		addedBy:      addedBy,
	}
	result.SetId(id)
	return result
}

func NewEnqueuedSong(
	songData song_data.SongData,
	addedBy user_id.UserId,
) *EnqueuedSong {
	result := &EnqueuedSong{
		addedAtUtc:   time.Now().UTC(),
		startedAtUtc: nil,
		addedBy:      addedBy,
		state:        enqueued_song_state.Enqueued,
		votes:        0,
		songData:     songData,
	}
	result.SetId(enqueued_song_id.NewEnqueuedSongId())
	return result
}

func (s EnqueuedSong) Equal(o EnqueuedSong) bool {
	if s.Id() != o.Id() {
		return false
	}
	if !s.songData.Equal(o.songData) {
		return false
	}
	if !nearlyEqual(s.addedAtUtc, o.addedAtUtc, time.Second) {
		return false
	}
	if (s.startedAtUtc == nil) != (o.startedAtUtc == nil) {
		return false
	}
	if s.startedAtUtc != nil && !nearlyEqual(*s.startedAtUtc, *o.startedAtUtc, time.Second) {
		return false
	}
	if s.state != o.state || s.votes != o.votes || s.addedBy != o.addedBy {
		return false
	}
	return true
}

func nearlyEqual(a, b time.Time, d time.Duration) bool {
	diff := a.Sub(b)
	if diff < 0 {
		diff = -diff
	}
	return diff <= d
}
