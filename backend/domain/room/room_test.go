package room

import (
	"fmt"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestNewRoomSuccess(t *testing.T) {
	hostID := user_id.NewUserId()
	name := faker.Name()
	password := gofakeit.Password(true, true, true, true, false, 15)
	qrCode := faker.UUIDDigit()
	encrypter := mock_encrypter.MockEncrypter{}
	encrypter.On("HashAndSalt", password).Return([]byte(password), nil)
	room, err := NewRoom(name, password, qrCode, hostID, &encrypter)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, name, room.Name())
	assert.Equal(t, []byte(password), room.Password())
	assert.Equal(t, qrCode, room.QrCode())
	assert.Len(t, room.Members(), 1)
	assert.Nil(t, room.BoostCooldownSeconds())
	assert.Equal(t, hostID, room.Members()[0])
	assert.WithinDuration(t, time.Now().UTC(), room.CreatedAtUtc(), time.Second)
	assert.Equal(t, uint32(DefaultLifespanSeconds), room.LifespanSeconds())
	assert.Zero(t, len(room.EnqueuedSongs()))
}

func TestNewRoomNameIncorrectLength(t *testing.T) {
	hostID := user_id.NewUserId()
	password := gofakeit.Password(true, true, true, true, false, 15)
	qrCode := faker.UUIDDigit()
	// Too short name
	var name string
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMinLength-1))
	encrypter := mock_encrypter.MockEncrypter{}
	encrypter.On("HashAndSalt", password).Return([]byte(password), nil)

	_, err := NewRoom(name, password, qrCode, hostID, &encrypter)
	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError := err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.Name.WrongLength")
	assert.Contains(t, domainError.Description, fmt.Sprintf("The room name has to be between %d and %d characters.",
		NameMinLength,
		NameMaxLength,
	))

	// Too long name
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMaxLength+1))
	_, err = NewRoom(name, password, qrCode, hostID, &encrypter)
	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError = err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.Name.WrongLength")
	assert.Contains(t, domainError.Description, fmt.Sprintf("The room name has to be between %d and %d characters.",
		NameMinLength,
		NameMaxLength,
	))
}

func TestNewRoomPasswordIncorrectLength(t *testing.T) {
	hostID := user_id.NewUserId()
	var name string
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMaxLength-1))
	qrCode := faker.UUIDDigit()

	// Too short password
	var password string
	encrypter := mock_encrypter.MockEncrypter{}
	encrypter.On("HashAndSalt", password).Return([]byte(password), nil)
	_ = faker.FakeData(&password, options.WithRandomStringLength(PasswordMinLength-1))
	_, err := NewRoom(name, password, qrCode, hostID, &encrypter)
	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError := err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.Password.WrongLength")
	assert.Contains(t, domainError.Description, fmt.Sprintf("The room password has to be between %d and %d characters.",
		PasswordMinLength,
		PasswordMaxLength,
	))

	// Too long password
	_ = faker.FakeData(&password, options.WithRandomStringLength(PasswordMaxLength+1))
	_, err = NewRoom(name, password, qrCode, hostID, &encrypter)
	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError = err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.Password.WrongLength")
	assert.Contains(t, domainError.Description, fmt.Sprintf("The room password has to be between %d and %d characters.",
		PasswordMinLength,
		PasswordMaxLength,
	))
}

func TestNewRoomQrCodeEmpty(t *testing.T) {
	hostID := user_id.NewUserId()
	qrCode := ""
	password := "validPassword"
	name := "Valid Room Name"

	encrypter := mock_encrypter.MockEncrypter{}
	encrypter.On("HashAndSalt", password).Return([]byte(password), nil)
	// Too short password
	_, err := NewRoom(name, password, qrCode, hostID, &encrypter)

	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError := err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.QrCode.Empty")
	assert.Equal(t, domainError.Description, "The room QR code cannot be empty.")
}
