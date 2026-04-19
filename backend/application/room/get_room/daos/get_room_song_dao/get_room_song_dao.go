package get_room_song_dao

import (
	"github.com/google/uuid"
)

type GetRoomSongDao struct {
	Id            uuid.UUID `db:"id"`
	Title         string    `db:"title"`
	Author        string    `db:"author"`
	AlbumCoverUrl string    `db:"album_cover_url"`
	State         string    `db:"state"`
	Votes         uint8     `db:"votes"`
	AddedBy       string    `db:"added_by"`
	VoteStatus    string    `db:"vote_status"`
}
