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
	author        string
	lengthSeconds uint16
	addedBy       user.Id
	addedAtUtc    time.Time
	startedAtUtc  *time.Time
	state         SongState
	votes         uint8
	wasUpVoted    bool
	wasDownVoted  bool
	albumCoverUrl string
}

func (s Song) AlbumCoverUrl() string {
	return s.albumCoverUrl
}

func (s Song) WasDownVoted() bool {
	return s.wasDownVoted
}

func (s Song) WasUpVoted() bool {
	return s.wasUpVoted
}

func (s Song) WasPlayed() bool {
	return s.State() == Played
}

func (s Song) ExternalId() string {
	return s.externalId
}

func (s Song) Title() string {
	return s.title
}

func (s Song) Author() string {
	return s.author
}

func (s Song) LengthSeconds() uint16 {
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

func NewSong(externalId, title, author, albumCoverUrl string, lengthSeconds uint16, addedBy user.Id) *Song {
	return &Song{
		externalId:    externalId,
		albumCoverUrl: albumCoverUrl,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
		addedBy:       addedBy,
		addedAtUtc:    time.Now().UTC(),
		startedAtUtc:  nil,
		state:         Enqueued,
		votes:         0,
		wasUpVoted:    false,
		wasDownVoted:  false,
	}
}
func HydrateSong(
	id SongId,
	externalId string,
	title string,
	author string,
	lengthSeconds uint16,
	addedBy user.Id,
	addedAtUtc time.Time,
	startedAtUtc *time.Time,
	state SongState,
	votes uint8,
	wasUpVoted bool,
	wasDownVoted bool,
	albumCoverUrl string,
) *Song {
	result := &Song{
		externalId:    externalId,
		title:         title,
		author:        author,
		lengthSeconds: lengthSeconds,
		addedBy:       addedBy,
		addedAtUtc:    addedAtUtc,
		startedAtUtc:  startedAtUtc,
		state:         state,
		votes:         votes,
		wasUpVoted:    wasUpVoted,
		wasDownVoted:  wasDownVoted,
		albumCoverUrl: albumCoverUrl,
	}
	result.SetId(id)
	return result
}
