package custom_errors

type Type int

const (
	Validation   Type = 400
	Unexpected   Type = 500
	NotFound     Type = 404
	Unauthorized Type = 401
	Forbidden    Type = 403
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
func (t Type) Error() string {
	return t.String()
}
