package user

type DeviceState int

const (
	Online DeviceState = iota
	Offline
)

var deviceStateName = map[DeviceState]string{
	Online:  "online",
	Offline: "offline",
}
var deviceStateFromName = map[string]DeviceState{
	"online":  Online,
	"offline": Offline,
}

func (s DeviceState) String() string {
	return deviceStateName[s]
}
func ParseDeviceState(s *string) *DeviceState {
	if s == nil {
		return nil
	}
	deviceState, ok := deviceStateFromName[*s]
	if !ok {
		return nil
	}
	return &deviceState
}
