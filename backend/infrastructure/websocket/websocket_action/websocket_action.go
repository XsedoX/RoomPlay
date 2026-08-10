package websocket_action

type (
	WebSocketAction         string
	WebSocketIncomingAction WebSocketAction
	WebSocketOutgoingAction WebSocketAction
)

const (
	SongEnqueuedClientMessageActionName = WebSocketIncomingAction("song_enqueued")

	RoomDataActionName           = WebSocketOutgoingAction("room_data")
	EnqueuedSongsPatchActionName = WebSocketOutgoingAction("enqueued_songs_patch")
	SongEnqueuedErrorActionName  = WebSocketOutgoingAction("song_enqueued_error")
)
