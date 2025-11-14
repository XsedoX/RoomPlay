package user

type UserRole int

const (
	Host UserRole = iota
	Member
)

var userRoleName = map[UserRole]string{
	Host:   "host",
	Member: "member",
}
var userRoleFromName = map[string]UserRole{
	"host":   Host,
	"member": Member,
}

func (r *UserRole) String() *string {
	if r == nil {
		return nil
	}
	s := userRoleName[*r]
	return &s
}
func ParseUserRole(s *string) *UserRole {
	if s == nil {
		return nil
	}
	role, ok := userRoleFromName[*s]
	if !ok {
		return nil
	}
	return &role
}
