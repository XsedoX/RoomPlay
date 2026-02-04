package room

import (
	"fmt"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
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

func (r Room) EnqueuedSongs() []enqueued_song.EnqueuedSong {
	return r.enqueuedSongs
}

func (r Room) Members() []user_id.UserId {
	return r.members
}

func NewRoom(name string,
	password string,
	qrCode string,
	roomHostId user_id.UserId,
) (*Room, error) {
	if len(name) > NameMaxLength {
		return nil, validation_domain_error.NewValidationDomainError("Room.TooLong.Name",
			fmt.Sprintf("The room name exceeded %d characters.",
				NameMaxLength))
	}
	if len(name) < NameMinLength {
		return nil, validation_domain_error.NewValidationDomainError("Room.TooShort.Name",
			fmt.Sprintf("The room name was shorter than %d characters.",
				NameMinLength))
	}
	if len(password) > PasswordMaxLength {
		return nil, validation_domain_error.NewValidationDomainError("Room.TooLong.Password",
			fmt.Sprintf("The room password exceeded %d characters.",
				PasswordMaxLength))
	}
	if len(password) < PasswordMinLength {
		return nil, validation_domain_error.NewValidationDomainError("Room.TooShort.Password",
			fmt.Sprintf("The room password was shorter than %d characters.",
				PasswordMinLength))
	}
	if qrCode == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"Room.QrCode",
			"qr code",
		)
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
