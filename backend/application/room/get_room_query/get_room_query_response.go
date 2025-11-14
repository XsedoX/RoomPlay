package get_room_query

import (
	"time"

	"github.com/google/uuid"
)

type GetRoomQueryResponse struct {
	Name        string            `json:"name"`
	BoostData   *BoostDataDto     `json:"boostData"`
	Songs       []RoomSongListDto `json:"songs"`
	QrCode      string            `json:"qrCode"`
	PlayingSong *PlayingSongDto   `json:"playingSong"`
	UserRole    string            `json:"userRole"`
}
type RoomSongListDto struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	AddedBy       string    `json:"addedBy"`
	Votes         uint8     `json:"votes"`
	AlbumCoverUrl string    `json:"albumCoverUrl"`
	Id            uuid.UUID `json:"id"`
	State         string    `json:"state"`
	VoteStatus    string    `json:"voteStatus"`
}
type PlayingSongDto struct {
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	StartedAtUtc  time.Time `json:"startedAtUtc"`
	LengthSeconds uint8     `json:"lengthSeconds"`
}
type BoostDataDto struct {
	BoostUsedAtUtc       time.Time `json:"boostUsedAtUtc"`
	BoostCooldownSeconds uint8     `json:"boostCooldownSeconds"`
}
