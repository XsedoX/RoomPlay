package application_error_type

type ApplicationErrorType int

const (
	Validation   ApplicationErrorType = 400
	Unexpected   ApplicationErrorType = 500
	NotFound     ApplicationErrorType = 404
	Unauthorized ApplicationErrorType = 401
	Forbidden    ApplicationErrorType = 403
)

var typeNames = map[ApplicationErrorType]string{
	Validation:   "validation",
	Unexpected:   "unexpected",
	NotFound:     "notFound",
	Unauthorized: "unauthorized",
	Forbidden:    "forbidden",
}

var errorFromName = map[string]ApplicationErrorType{
	"validation":   Validation,
	"unexpected":   Unexpected,
	"notFound":     NotFound,
	"unauthorized": Unauthorized,
	"forbidden":    Forbidden,
}

func (t ApplicationErrorType) String() string {
	return typeNames[t]
}

func ParseType(s string) (ApplicationErrorType, bool) {
	v, ok := errorFromName[s]
	return v, ok
}

func (t ApplicationErrorType) Error() string {
	return t.String()
}
