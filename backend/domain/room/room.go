package room

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/domainErrors"
	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

const (
	NameMaxLength          = 30
	NameMinLength          = 5
	PasswordMaxLength      = 30
	PasswordMinLength      = 10
	DefaultLifespanSeconds = 60 * 60 * 24 // 24 hours
)

type Room struct {
	shared.AggregateRoot[shared.RoomId]
	name                 string
	password             string
	qrCode               string
	boostCooldownSeconds *uint8
	createdAtUtc         time.Time
	lifespanSeconds      uint32
	songsList            []Song
	members              []user.Id
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

func (r Room) BoostCooldownSeconds() *uint8 {
	return r.boostCooldownSeconds
}

func (r Room) CreatedAtUtc() time.Time {
	return r.createdAtUtc
}

func (r Room) LifespanSeconds() uint32 {
	return r.lifespanSeconds
}

func (r Room) SongsList() []Song {
	return r.songsList
}

func (r Room) Members() []user.Id {
	return r.members
}

func NewRoom(name string,
	password string,
	qrCode string,
	roomHostId user.Id) *Room {

	if len(name) > NameMaxLength {
		panic(domainErrors.NewValidationDomainError("Room.TooLong.Name",
			fmt.Sprintf("The room name exceeded %d characters.",
				NameMaxLength)))
	}
	if len(name) < NameMinLength {
		panic(domainErrors.NewValidationDomainError("Room.TooShort.Name",
			fmt.Sprintf("The room was shorter than %d characters.",
				NameMinLength)))
	}
	if len(password) > PasswordMaxLength {
		panic(domainErrors.NewValidationDomainError("Room.TooLong.Password",
			fmt.Sprintf("The room password exceeded %d characters.",
				PasswordMaxLength)))
	}
	if len(password) < PasswordMinLength {
		panic(domainErrors.NewValidationDomainError("Room.TooShort.Password",
			fmt.Sprintf("The room password was shorter than %d characters.",
				PasswordMinLength)))
	}
	result := &Room{
		name:                 name,
		password:             password,
		qrCode:               qrCode,
		boostCooldownSeconds: nil,
		createdAtUtc:         time.Now().UTC(),
		lifespanSeconds:      DefaultLifespanSeconds,
		songsList:            make([]Song, 0),
		members:              []user.Id{roomHostId},
	}
	result.SetId(shared.RoomId(uuid.New()))
	return result
}
