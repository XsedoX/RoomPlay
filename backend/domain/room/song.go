package room

import (
	"time"

	"xsedox.com/main/domain/shared"
	"xsedox.com/main/domain/user"
)

type Song struct {
	shared.Entity[SongId]
	externalId    string
	title         string
	artist        string
	lengthSeconds uint8
	addedBy       user.Id
	addedAtUtc    time.Time
	startedAtUtc  *time.Time
	state         SongState
	votes         uint8
}

func (s Song) ExternalId() string {
	return s.externalId
}

func (s Song) Title() string {
	return s.title
}

func (s Song) Artist() string {
	return s.artist
}

func (s Song) LengthSeconds() uint8 {
	return s.lengthSeconds
}

func (s Song) AddedBy() user.Id {
	return s.addedBy
}

func (s Song) AddedAtUtc() time.Time {
	return s.addedAtUtc
}

func (s Song) StartedAtUtc() *time.Time {
	return s.startedAtUtc
}

func (s Song) State() SongState {
	return s.state
}

func (s Song) Votes() uint8 {
	return s.votes
}

func NewSong(externalId string, title string, artist string, lengthSeconds uint8, addedBy user.Id) *Song {
	return &Song{
		externalId:    externalId,
		title:         title,
		artist:        artist,
		lengthSeconds: lengthSeconds,
		addedBy:       addedBy,
		addedAtUtc:    time.Now().UTC(),
		startedAtUtc:  nil,
		state:         Enqueued,
		votes:         0,
	}
}
