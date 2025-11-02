package device

type State int

const (
	Online State = iota
	Offline
)

var stateName = map[State]string{
	Online:  "online",
	Offline: "offline",
}
var stateFromName = map[string]State{
	"online":  Online,
	"offline": Offline,
}

func (s State) String() string {
	return stateName[s]
}
func ParseState(s string) State {
	deviceState, ok := stateFromName[s]
	if !ok {
		panic("Unknown device state: " + s)
	}
	return deviceState
}
