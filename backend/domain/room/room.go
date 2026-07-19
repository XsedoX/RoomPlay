package room

import (
	"bytes"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/room/default_playlist"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/room/events"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/room/scheduled_song"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/google/uuid"
)

const (
	NameMaxLength          = 30
	NameMinLength          = 5
	PasswordMaxLength      = 30
	PasswordMinLength      = 10
	DefaultLifespanSeconds = 60 * 60 * 24 // 24 hours
)

type Room struct {
	shared.AggregateRoot[room_id.RoomId]
	name                 string
	password             []byte
	qrCode               string
	boostCooldownSeconds *uint16
	createdAtUtc         time.Time
	lifespanSeconds      uint32
	enqueuedSongs        []enqueued_song.EnqueuedSong
	members              []user_id.UserId
	bannedUsers          []user_id.UserId
	scheduledSong        *scheduled_song.ScheduledSong
	defaultPlaylist      *default_playlist.DefaultPlaylist
}

func (r Room) ScheduledSong() *scheduled_song.ScheduledSong {
	return r.scheduledSong
}

func (r Room) DefaultPlaylist() *default_playlist.DefaultPlaylist {
	return r.defaultPlaylist
}

func (r *Room) EnqueueSong(addedBy user_id.UserId, songData song_data.SongData) {
	enqueuedSong := enqueued_song.NewEnqueuedSong(songData, addedBy)
	r.enqueuedSongs = append(r.enqueuedSongs, *enqueuedSong)
	songEnqueuedEvent := events.NewSongEnqueuedEvent(
		enqueuedSong.Id(),
		songData.Title(),
		songData.Author(),
		addedBy,
		enqueuedSong.Votes(),
		songData.AlbumCoverUrl(),
		enqueuedSong.State(),
		vote_status.NotVoted,
	)
	r.RaiseDomainEvent(songEnqueuedEvent)
}

func (r Room) PlayingSong() *enqueued_song.EnqueuedSong {
	for _, song := range r.enqueuedSongs {
		if song.State() == enqueued_song_state.Playing {
			return &song
		}
	}
	return nil
}

func (r Room) Name() string {
	return r.name
}

func (r Room) Password() []byte {
	return r.password
}

func (r Room) QrCode() string {
	return r.qrCode
}

func (r Room) BoostCooldownSeconds() *uint16 {
	return r.boostCooldownSeconds
}

func (r Room) CreatedAtUtc() time.Time {
	return r.createdAtUtc
}

func (r Room) LifespanSeconds() uint32 {
	return r.lifespanSeconds
}

func (r Room) AllSongs() []enqueued_song.EnqueuedSong {
	return r.enqueuedSongs
}

func (r Room) EnqueuedSongs() []enqueued_song.EnqueuedSong {
	result := make([]enqueued_song.EnqueuedSong, 0, len(r.enqueuedSongs))
	for _, song := range r.enqueuedSongs {
		if song.State() != enqueued_song_state.Playing {
			result = append(result, song)
		}
	}
	return result
}

func (r Room) Members() []user_id.UserId {
	return r.members
}

func (r Room) BannedUsers() []user_id.UserId {
	return r.bannedUsers
}

func NewRoom(name string,
	password string,
	qrCode string,
	roomHostId user_id.UserId,
	hasher IPasswordHasher,
) (*Room, error) {
	if (len(name) > NameMaxLength) ||
		(len(name) < NameMinLength) {
		return nil, domain_errors.NewRoomNameIncorrectError(
			NameMaxLength,
			NameMinLength,
		)
	}
	if (len(password) > PasswordMaxLength) ||
		(len(password) < PasswordMinLength) {
		return nil, domain_errors.NewRoomPasswordIncorrectError(
			PasswordMaxLength,
			PasswordMinLength,
		)
	}
	if qrCode == "" {
		return nil, domain_errors.NewRoomQrCodeEmptyError()
	}
	passwdordBytes, hashErr := hasher.HashAndSalt(password)
	if hashErr != nil {
		return nil, domain_errors.NewRoomHashSaltError()
	}
	result := &Room{
		name:                 name,
		password:             passwdordBytes,
		qrCode:               qrCode,
		boostCooldownSeconds: nil,
		createdAtUtc:         time.Now().UTC(),
		lifespanSeconds:      DefaultLifespanSeconds,
		enqueuedSongs:        make([]enqueued_song.EnqueuedSong, 0),
		members:              []user_id.UserId{roomHostId},
	}
	result.SetId(room_id.RoomId(uuid.New()))
	return result, nil
}

func HydrateRoom(
	id room_id.RoomId,
	name string,
	password []byte,
	qrCode string,
	boostCooldownSeconds *uint16,
	createdAtUtc time.Time,
	lifespanSeconds uint32,
	enqueuedSongs []enqueued_song.EnqueuedSong,
	members []user_id.UserId,
) *Room {
	r := &Room{
		name:                 name,
		password:             password,
		qrCode:               qrCode,
		boostCooldownSeconds: boostCooldownSeconds,
		createdAtUtc:         createdAtUtc,
		lifespanSeconds:      lifespanSeconds,
		enqueuedSongs:        enqueuedSongs,
		members:              members,
	}
	r.SetId(id)
	return r
}

func (r Room) Equal(o Room) bool {
	if r.Id() != o.Id() {
		return false
	}
	if r.name != o.name || !bytes.Equal(r.Password(), o.Password()) || r.qrCode != o.qrCode {
		return false
	}
	if (r.boostCooldownSeconds == nil) != (o.boostCooldownSeconds == nil) {
		return false
	}
	if r.boostCooldownSeconds != nil && *r.boostCooldownSeconds != *o.boostCooldownSeconds {
		return false
	}
	if !nearlyEqual(r.createdAtUtc, o.createdAtUtc, time.Second) {
		return false
	}
	if r.lifespanSeconds != o.lifespanSeconds {
		return false
	}
	if len(r.enqueuedSongs) != len(o.enqueuedSongs) {
		return false
	}
	for i := range r.enqueuedSongs {
		if !r.enqueuedSongs[i].Equal(o.enqueuedSongs[i]) {
			return false
		}
	}
	if len(r.members) != len(o.members) {
		return false
	}
	if !usersEqual(r.members, o.members) {
		return false
	}
	if len(r.bannedUsers) != len(o.bannedUsers) {
		return false
	}
	if !usersEqual(r.bannedUsers, o.bannedUsers) {
		return false
	}
	if (r.scheduledSong == nil) != (o.scheduledSong == nil) {
		return false
	}
	if r.scheduledSong != nil && *r.scheduledSong != *o.scheduledSong {
		return false
	}
	if (r.defaultPlaylist == nil) != (o.defaultPlaylist == nil) {
		return false
	}
	if r.defaultPlaylist != nil && *r.defaultPlaylist != *o.defaultPlaylist {
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

// ...existing code...
func usersEqual(a, b []user_id.UserId) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[user_id.UserId]struct{}, len(a))
	for _, u := range a {
		seen[u] = struct{}{}
	}
	for _, u := range b {
		if _, ok := seen[u]; !ok {
			return false
		}
		delete(seen, u)
	}
	return len(seen) == 0
}
