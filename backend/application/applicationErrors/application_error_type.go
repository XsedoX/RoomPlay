package applicationErrors

type Type int

const (
	Validation Type = iota
	Unexpected
	NotFound
	Unauthorized
	Forbidden
)

var typeNames = map[Type]string{
	Validation:   "Validation",
	Unexpected:   "Unexpected",
	NotFound:     "NotFound",
	Unauthorized: "Unauthorized",
	Forbidden:    "Forbidden",
}
var errorFromName = map[string]Type{
	"Validation":   Validation,
	"Unexpected":   Unexpected,
	"NotFound":     NotFound,
	"Unauthorized": Unauthorized,
	"Forbidden":    Forbidden,
}

func (t Type) String() string {
	return typeNames[t]
}
func ParseType(s string) (Type, bool) {
	v, ok := errorFromName[s]
	return v, ok
}
