package user

type Role int

const (
	Host Role = iota
	Guest
)

var roleName = map[Role]string{
	Host:  "host",
	Guest: "guest",
}
var roleFromName = map[string]Role{
	"host":  Host,
	"guest": Guest,
}

func (r Role) String() string {
	return roleName[r]
}
func ParseRole(s string) (Role, bool) {
	role, ok := roleFromName[s]
	return role, ok
}
