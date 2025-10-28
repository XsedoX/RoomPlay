package contracts

type IEncrypter interface {
	Encrypt(plaintext string) (ciphertext []byte, err error)
	Decrypt(ciphertext []byte) (plaintext string, err error)
	NewEncryptionKey() []byte
	HashAndSalt(plaintext string) (hash []byte, err error)
	Verify(plaintext string, hash []byte) (ok bool)
}
