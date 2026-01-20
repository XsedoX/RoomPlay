package room

import (
	"testing"
	"time"

	domainErrors "github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/user"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRoomSuccess(t *testing.T) {
	hostID := user.Id(uuid.New())
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
	assert.Zero(t, len(room.SongsList()))
}

func TestNewRoomNameIncorrectLength(t *testing.T) {
	hostID := user.Id(uuid.New())
	var password string
	_ = faker.FakeData(&password, options.WithRandomStringLength(20))
	qrCode := faker.UUIDDigit()
	// Too short name
	var name string
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMinLength-1))

	_, err := NewRoom(name, password, qrCode, hostID)
	assert.Error(t, err)
	assert.IsType(t, &domainErrors.ValidationDomainError{}, err)
	domainError := err.(*domainErrors.ValidationDomainError)
	assert.Equal(t, domainError.Code, "Room.TooShort.Name")
	assert.Contains(t, domainError.Description, "The room name was shorter than")
	assert.Equal(t, domainError.Title, "Validation error occurred.")

	// Too long name
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMaxLength+1))
	_, err = NewRoom(name, password, qrCode, hostID)
	assert.Error(t, err)
	assert.IsType(t, &domainErrors.ValidationDomainError{}, err)
	domainError = err.(*domainErrors.ValidationDomainError)
	assert.Equal(t, domainError.Code, "Room.TooLong.Name")
	assert.Contains(t, domainError.Description, "The room name exceeded")
	assert.Equal(t, domainError.Title, "Validation error occurred.")
}

func TestNewRoomPasswordIncorrectLength(t *testing.T) {
	hostID := user.Id(uuid.New())
	var name string
	_ = faker.FakeData(&name, options.WithRandomStringLength(NameMaxLength-1))
	qrCode := faker.UUIDDigit()

	// Too short password
	var password string
	_ = faker.FakeData(&password, options.WithRandomStringLength(PasswordMinLength-1))
	_, err := NewRoom(name, password, qrCode, hostID)
	assert.Error(t, err)
	assert.IsType(t, &domainErrors.ValidationDomainError{}, err)
	domainError := err.(*domainErrors.ValidationDomainError)
	assert.Equal(t, domainError.Code, "Room.TooShort.Password")
	assert.Equal(t, domainError.Title, "Validation error occurred.")
	assert.Contains(t, domainError.Description, "The room password was shorter than")

	// Too long password
	_ = faker.FakeData(&password, options.WithRandomStringLength(PasswordMaxLength+1))
	_, err = NewRoom(name, password, qrCode, hostID)
	assert.Error(t, err)
	assert.IsType(t, &domainErrors.ValidationDomainError{}, err)
	domainError = err.(*domainErrors.ValidationDomainError)
	assert.Equal(t, domainError.Code, "Room.TooLong.Password")
	assert.Contains(t, domainError.Description, "The room password exceeded")
	assert.Equal(t, domainError.Title, "Validation error occurred.")
}
