package get_room_query

import (
	"time"

	"github.com/google/uuid"
)

type GetRoomQueryResponse struct {
	Name           string            `json:"name"`
	Songs          []RoomSongListDto `json:"songs"`
	QrCode         string            `json:"qrCode"`
	BoostUsedAtUtc *time.Time        `json:"boostUsedAtUtc"`
	PlayingSong    *PlayingSongDto   `json:"playingSong"`
	UserRole       string            `json:"userRole"`
}
type RoomSongListDto struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	AddedBy       string    `json:"addedBy"`
	Votes         uint8     `json:"votes"`
	AlbumCoverUrl string    `json:"albumCoverUrl"`
	Id            uuid.UUID `json:"id"`
	State         string    `json:"state"`
}
type PlayingSongDto struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	StartedAtUtc  time.Time `json:"startedAtUtc"`
	LengthSeconds uint8     `json:"lengthSeconds"`
}
