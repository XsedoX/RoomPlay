package events

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/shared"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

const SongEnqueuedEventName shared.EventName = "SongEnqueuedEvent"

type SongEnqueuedEvent struct {
	enqueuedSongId    enqueued_song_id.EnqueuedSongId
	title             string
	author            string
	addedBy           user_id.UserId
	votes             int8
	albumCoverUrl     string
	enqueuedSongState enqueued_song_state.EnqueuedSongState
	status            vote_status.VoteStatus
	roomId            room_id.RoomId

	eventName  shared.EventName
	occurredAt time.Time
}

func (event *SongEnqueuedEvent) EnqueuedSongId() enqueued_song_id.EnqueuedSongId {
	return event.enqueuedSongId
}

func (event *SongEnqueuedEvent) RoomId() room_id.RoomId {
	return event.roomId
}

func (event *SongEnqueuedEvent) Title() string {
	return event.title
}

func (event *SongEnqueuedEvent) Author() string {
	return event.author
}

func (event *SongEnqueuedEvent) AddedBy() user_id.UserId {
	return event.addedBy
}

func (event *SongEnqueuedEvent) Votes() int8 {
	return event.votes
}

func (event *SongEnqueuedEvent) AlbumCoverUrl() string {
	return event.albumCoverUrl
}

func (event *SongEnqueuedEvent) EnqueuedSongState() enqueued_song_state.EnqueuedSongState {
	return event.enqueuedSongState
}

func (event *SongEnqueuedEvent) Status() vote_status.VoteStatus {
	return event.status
}

func NewSongEnqueuedEvent(
	enqueuedSongId enqueued_song_id.EnqueuedSongId,
	title string,
	author string,
	addedBy user_id.UserId,
	votes int8,
	albumCoverUrl string,
	enqueuedSongState enqueued_song_state.EnqueuedSongState,
	status vote_status.VoteStatus,
) *SongEnqueuedEvent {
	return &SongEnqueuedEvent{
		enqueuedSongId:    enqueuedSongId,
		title:             title,
		author:            author,
		addedBy:           addedBy,
		votes:             votes,
		albumCoverUrl:     albumCoverUrl,
		enqueuedSongState: enqueuedSongState,
		status:            status,

		eventName:  SongEnqueuedEventName,
		occurredAt: time.Now(),
	}
}

func (event *SongEnqueuedEvent) EventName() shared.EventName {
	return event.eventName
}

func (event *SongEnqueuedEvent) OccurredAt() time.Time {
	return event.occurredAt
}
