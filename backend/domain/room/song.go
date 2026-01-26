package room

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user"
)

type Song struct {
	shared.Entity[SongId]
	url           string
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

func (s Song) Url() string {
	return s.url
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

func NewSong(url, title, author, albumCoverUrl string, lengthSeconds uint16, addedBy user.Id) *Song {
	return &Song{
		url:           url,
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
	url string,
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
		url:           url,
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
