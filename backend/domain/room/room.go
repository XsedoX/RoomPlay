package room

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/room/default_playlist"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
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
	password             string
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

func (r Room) Password() string {
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
	result := &Room{
		name:                 name,
		password:             password,
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
	password string,
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
