package room

type SongState int

const (
	Enqueued SongState = iota
	Playing
	Played
)

var songStateName = map[SongState]string{
	Enqueued: "enqueued",
	Playing:  "playing",
	Played:   "played",
}
var songStateValue = map[string]SongState{
	"enqueued": Enqueued,
	"playing":  Playing,
	"played":   Played,
}

func (s SongState) String() string {
	return songStateName[s]
}
func ParseSongState(s string) *SongState {
	songState, ok := songStateValue[s]
	if !ok {
		return nil
	}
	return &songState
}
