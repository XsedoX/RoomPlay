package default_playlist_dao

import (
	"github.com/XsedoX/RoomPlay/domain/room/default_playlist"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/google/uuid"
)

type DefaultPlaylistDao struct {
	RoomId        uuid.UUID        `gorm:"primaryKey;foreignKey:RoomId;references:Id"`
	ExternalId    string           `gorm:"column:external_id"`
	UserId        uuid.UUID        `gorm:"column:user_id;foreignKey:UserId;references:Id"`
	UserDao       user_dao.UserDao `gorm:"foreignKey:UserId;references:Id"`
	SongsAmount   uint             `gorm:"column:songs_amount;type:smallint"`
	PlaylistTitle string           `gorm:"column:playlist_title"`
}

func (DefaultPlaylistDao) TableName() string {
	return "default_playlists"
}

func (d DefaultPlaylistDao) ToDomain() *default_playlist.DefaultPlaylist {
	return default_playlist.HydrateDefaultPlaylist(
		d.ExternalId,
		user_id.UserId(d.UserId),
		d.SongsAmount,
		d.PlaylistTitle,
	)
}

func (d *DefaultPlaylistDao) FromDomain(value *default_playlist.DefaultPlaylist, roomId room_id.RoomId) {
	d.RoomId = roomId.ToUuid()
	d.ExternalId = value.ExternalId()
	d.UserId = value.UserId().ToUuid()
	d.SongsAmount = value.SongsAmount()
	d.PlaylistTitle = value.Title()
	// The owning user is a separate aggregate, not part of DefaultPlaylist.
	d.UserDao = user_dao.UserDao{}
}
