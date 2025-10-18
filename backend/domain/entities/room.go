package entities

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const DefaultLifespanSeconds = 60 * 60 * 24 // 24 hours

type RoomId uuid.UUID
type Room struct {
	Entity[RoomId]
	roomName             string
	roomPassword         string
	qrCode               string
	boostCooldownSeconds *int
	createdAtUtc         time.Time
	lifespanSeconds      int
	userIds              []UserId
}

func NewRoom(roomName string, roomPassword string, userId UserId) (*Room, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		// If reading random bytes fails, return the error.
		return nil, err
	}
	// Encode the random bytes into a hex string.
	qrCodeHash := hex.EncodeToString(randomBytes)

	room := &Room{
		roomName:             roomName,
		roomPassword:         roomPassword,
		createdAtUtc:         time.Now().UTC(),
		boostCooldownSeconds: nil, // Default to nil
		qrCode:               qrCodeHash,
		userIds:              []UserId{userId},
		lifespanSeconds:      DefaultLifespanSeconds, // Default lifespan of 24 hours
	}
	room.SetId(RoomId(uuid.New()))
	return room, nil
}

func (room *Room) GetId() RoomId {
	return room.GetId()
}
func (room *Room) GetRoomName() string {
	return room.roomName
}
func (room *Room) GetRoomPassword() string {
	return room.roomPassword
}
func (room *Room) GetQrCode() string {
	return room.qrCode
}
func (room *Room) GetBoostCooldownSeconds() *int {
	return room.boostCooldownSeconds
}
func (room *Room) GetCreatedAtUtc() time.Time {
	return room.createdAtUtc
}
func (room *Room) GetLifespanSeconds() int {
	return room.lifespanSeconds
}
func (room *Room) GetUserIds() []UserId {
	return room.userIds
}
