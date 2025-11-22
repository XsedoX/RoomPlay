package authentication_mocks

import "github.com/stretchr/testify/mock"

type MockEncrypter struct {
	mock.Mock
}

func (m *MockEncrypter) Encrypt(plaintext string) (ciphertext []byte, err error) {
	args := m.Called(plaintext)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncrypter) Decrypt(ciphertext []byte) (plaintext string, err error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}

func (m *MockEncrypter) NewEncryptionKey() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *MockEncrypter) HashAndSalt(plaintext string) (hash []byte, err error) {
	args := m.Called(plaintext)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncrypter) Verify(plaintext string, hash []byte) (ok bool) {
	args := m.Called(plaintext, hash)
	return args.Bool(0)
}

func (m *MockEncrypter) Hash(plaintext string) []byte {
	args := m.Called(plaintext)
	return args.Get(0).([]byte)
}
