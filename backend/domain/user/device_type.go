package user

type DeviceType int

const (
	Mobile DeviceType = iota
	Desktop
)

var deviceTypeName = map[DeviceType]string{
	Mobile:  "mobile",
	Desktop: "desktop",
}
var deviceTypeFromName = map[string]DeviceType{
	"mobile":  Mobile,
	"desktop": Desktop,
}

func (t DeviceType) String() string {
	return deviceTypeName[t]
}
func ParseDeviceType(s *string) *DeviceType {
	if s == nil {
		return nil
	}
	deviceType, ok := deviceTypeFromName[*s]
	if !ok {
		return nil
	}
	return &deviceType
}
func ListDeviceTypes() []string {
	types := make([]string, 0, len(deviceTypeName))
	for _, deviceType := range deviceTypeName {
		types = append(types, deviceType)
	}
	return types
}
