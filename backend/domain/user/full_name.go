package user

import "fmt"

type FullName struct {
	name    string
	surname string
}

func NewFullName(name string, surname string) FullName {
	return FullName{
		name:    name,
		surname: surname,
	}
}

func (fn FullName) Equals(fnOther FullName) bool {
	return fn.name == fnOther.Name() && fn.surname == fnOther.Surname()
}
func (fn FullName) Name() string {
	return fn.name
}
func (fn FullName) Surname() string {
	return fn.surname
}
func (fn FullName) String() string {
	return fmt.Sprintf("%s %s", fn.name, fn.surname)
}
