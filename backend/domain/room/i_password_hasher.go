package room

type IPasswordHasher interface {
	HashAndSalt(password string) (hash []byte, err error)
}
