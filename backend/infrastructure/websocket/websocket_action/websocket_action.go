package websocket_action

type (
	WebSocketAction         string
	WebSocketIncomingAction WebSocketAction
	WebSocketOutgoingAction WebSocketAction
)

const (
	JoinedRoomActionName      = WebSocketOutgoingAction("joined_room")
	JoinedRoomErrorActionName = WebSocketOutgoingAction("joined_room_error")

	EnqueueSongClientMessageActionName = WebSocketIncomingAction("enqueue_song")
	EnqueuedSongPatchActionName        = WebSocketOutgoingAction("enqueued_song_patch")
	SongEnqueuedErrorActionName        = WebSocketOutgoingAction("song_enqueued_error")
)
