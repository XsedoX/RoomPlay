package device

type Type int

const (
	Mobile Type = iota
	Computer
)

var typeName = map[Type]string{
	Mobile:   "mobile",
	Computer: "computer",
}
var typeFromName = map[string]Type{
	"mobile":   Mobile,
	"computer": Computer,
}

func (t Type) String() string {
	return typeName[t]
}
func ParseType(s string) (Type, bool) {
	deviceType, ok := typeFromName[s]
	return deviceType, ok
}
