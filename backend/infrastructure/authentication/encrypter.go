// cryptopasta - basic cryptography examples
//
// Written in 2015 by George Tankersley <george.tankersley@gmail.com>
//
// To the extent possible under law, the author(s) have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// You should have received a copy of the CC0 Public Domain Dedication along
// with this software. If not, see // <http://creativecommons.org/publicdomain/zero/1.0/>.

// Provides symmetric authenticated encryption using 256-bit AES-GCM with a random nonce.
package authentication

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
	"xsedox.com/main/config"
)

type Encrypter struct {
	key [32]byte
}

func NewEncrypter(configuration config.IConfiguration) *Encrypter {
	var arr [32]byte
	copy(arr[:], configuration.Authentication().EncryptionKey)
	return &Encrypter{
		key: arr,
	}
}

// NewEncryptionKey generates a random 256-bit key for Encrypt() and
// Decrypt(). It panics if the source of randomness fails.
func (enc *Encrypter) NewEncryptionKey() []byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key[:]
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (enc *Encrypter) Encrypt(plaintextString string) (ciphertext []byte, err error) {
	plaintext := []byte(plaintextString)
	block, err := aes.NewCipher(enc.key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (enc *Encrypter) Decrypt(ciphertext []byte) (plaintext string, err error) {
	block, err := aes.NewCipher(enc.key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	bytePlaintext, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	return string(bytePlaintext), nil
}

func (enc *Encrypter) HashAndSalt(plaintext string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
func (enc *Encrypter) Verify(plaintext string, hash []byte) (ok bool) {
	ok = bcrypt.CompareHashAndPassword(hash, []byte(plaintext)) == nil
	return
}
