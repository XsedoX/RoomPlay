package daos

import (
	"github.com/google/uuid"
)

type GetRoomSongDao struct {
	Id            uuid.UUID `db:"id"`
	Title         string    `db:"title"`
	Author        string    `db:"author"`
	AddedBy       string    `db:"added_by"`
	State         string    `db:"state"`
	Votes         uint8     `db:"votes"`
	AlbumCoverUrl string    `db:"album_cover_url"`
	VoteStatus    string    `db:"vote_status"`
}
