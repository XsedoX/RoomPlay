package room

import (
	"fmt"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/stretchr/testify/assert"
)

func TestNewRoomSuccess(t *testing.T) {
	hostID := user_id.NewUserId()
	name := faker.Name()
	var password string
	_ = faker.FakeData(&password, options.WithRandomStringLength(20))
	qrCode := faker.UUIDDigit()

	room, err := NewRoom(name, password, qrCode, hostID)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, name, room.Name())
	assert.Equal(t, password, room.Password())
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
	var password string
	_ = faker.FakeData(&password, options.WithRandomStringLength(20))
	qrCode := faker.UUIDDigit()
	// Too short name
	var name string
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMinLength-1))

	_, err := NewRoom(name, password, qrCode, hostID)
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
	_, err = NewRoom(name, password, qrCode, hostID)
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
	_ = faker.FakeData(&password, options.WithRandomStringLength(PasswordMinLength-1))
	_, err := NewRoom(name, password, qrCode, hostID)
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
	_, err = NewRoom(name, password, qrCode, hostID)
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

	// Too short password
	_, err := NewRoom(name, password, qrCode, hostID)

	assert.Error(t, err)
	assert.IsType(t, &domain_errors.DomainError{}, err)
	domainError := err.(*domain_errors.DomainError)
	assert.Equal(t, domainError.Code, "Room.QrCode.Empty")
	assert.Equal(t, domainError.Description, "The room QR code cannot be empty.")
}
