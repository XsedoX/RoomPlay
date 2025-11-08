package applicationErrors

type Type int

const (
	Validation   Type = 400
	Unexpected        = 500
	NotFound          = 404
	Unauthorized      = 401
	Forbidden         = 403
)

var typeNames = map[Type]string{
	Validation:   "validation",
	Unexpected:   "unexpected",
	NotFound:     "notFound",
	Unauthorized: "unauthorized",
	Forbidden:    "forbidden",
}
var errorFromName = map[string]Type{
	"validation":   Validation,
	"unexpected":   Unexpected,
	"notFound":     NotFound,
	"unauthorized": Unauthorized,
	"forbidden":    Forbidden,
}

func (t Type) String() string {
	return typeNames[t]
}
func ParseType(s string) (Type, bool) {
	v, ok := errorFromName[s]
	return v, ok
}
