package enqueued_song_state

type EnqueuedSongState int

const (
	Enqueued EnqueuedSongState = iota
	Playing  EnqueuedSongState = iota
	Played   EnqueuedSongState = iota
	Paused   EnqueuedSongState = iota
)

var songStateName = map[EnqueuedSongState]string{
	Enqueued: "enqueued",
	Playing:  "playing",
	Played:   "played",
	Paused:   "paused",
}

var songStateValue = map[string]EnqueuedSongState{
	"enqueued": Enqueued,
	"playing":  Playing,
	"played":   Played,
	"paused":   Paused,
}

func (s EnqueuedSongState) String() string {
	return songStateName[s]
}

func ParseSongState(s string) *EnqueuedSongState {
	songState, ok := songStateValue[s]
	if !ok {
		return nil
	}
	return &songState
}
