package room

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/shared"
)

const DefaultLifespanSeconds = 60 * 60 * 24 // 24 hours

type Room struct {
	shared.Entity[shared.RoomId]
	roomName             string
	roomPassword         string
	qrCode               string
	boostCooldownSeconds *int
	createdAtUtc         time.Time
	lifespanSeconds      int
	userIds              []shared.UserId
}

func CreateRoom(roomName string, roomPassword string, userId shared.UserId) (*Room, error) {
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
		userIds:              []shared.UserId{userId},
		lifespanSeconds:      DefaultLifespanSeconds, // Default lifespan of 24 hours
	}
	room.SetId(shared.RoomId(uuid.New()))
	return room, nil
}

func (room *Room) GetId() shared.RoomId {
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
func (room *Room) GetUserIds() []shared.UserId {
	return room.userIds
}
