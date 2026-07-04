package enqueued_song_state

type EnqueuedSongState int

const (
	Enqueued EnqueuedSongState = iota
	Playing  EnqueuedSongState = iota
	Played   EnqueuedSongState = iota
)

var songStateName = map[EnqueuedSongState]string{
	Enqueued: "enqueued",
	Playing:  "playing",
	Played:   "played",
}

var songStateValue = map[string]EnqueuedSongState{
	"enqueued": Enqueued,
	"playing":  Playing,
	"played":   Played,
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
