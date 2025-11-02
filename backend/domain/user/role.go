package user

type Role int

const (
	Host Role = iota
	Member
)

var roleName = map[Role]string{
	Host:   "host",
	Member: "member",
}
var roleFromName = map[string]Role{
	"host":   Host,
	"member": Member,
}

func (r *Role) String() *string {
	if r == nil {
		return nil
	}
	s := roleName[*r]
	return &s
}
func ParseRole(s *string) *Role {
	if s == nil {
		return nil
	}
	role, ok := roleFromName[*s]
	if !ok {
		return nil
	}
	return &role
}
