package song_enqueued_websocket_event

import "github.com/google/uuid"

type WebSocketAction string

type SongEnqueuedWebsocketEventResponse struct {
	Title         string          `json:"title"`
	Author        string          `json:"author"`
	AddedBy       string          `json:"addedBy"`
	Votes         int8            `json:"votes"`
	AlbumCoverUrl string          `json:"albumCoverUrl"`
	Id            uuid.UUID       `json:"id"`
	State         string          `json:"state"`
	VoteStatus    string          `json:"voteStatus"`
	Action        WebSocketAction `json:"action"`
}
