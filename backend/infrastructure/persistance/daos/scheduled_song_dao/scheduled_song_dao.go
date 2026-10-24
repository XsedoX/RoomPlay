package scheduled_song_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/room/scheduled_song"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/song_dao"
	"github.com/google/uuid"
)

type ScheduledSongDao struct {
	RoomId         uuid.UUID        `gorm:"column:room_id;type:uuid;primaryKey;foreignKey:RoomId;references:Id" db:"room_id"`
	SongId         uuid.UUID        `gorm:"column:song_id;type:uuid;primaryKey;foreignKey:SongId;references:Id" db:"song_id"`
	SongDao        song_dao.SongDao `gorm:"foreignKey:SongId;references:Id" db:"song_data"`
	ScheduledAtUtc time.Time        `gorm:"column:scheduled_at_utc;type:timestamp" db:"scheduled_at_utc"`
}

func (ScheduledSongDao) TableName() string {
	return "scheduled_songs"
}

func (d ScheduledSongDao) ToDomain(provider music_provider.MusicProvider) *scheduled_song.ScheduledSong {
	return scheduled_song.HydrateScheduledSong(
		*d.SongDao.ToDomain(provider),
		d.ScheduledAtUtc,
	)
}

// FromDomain maps into the selected song DAO. When replacing the underlying song,
// supply its resolved SongDao/SongId, or use a fresh DAO for a new song record.
func (d *ScheduledSongDao) FromDomain(value *scheduled_song.ScheduledSong, roomId room_id.RoomId) {
	d.RoomId = roomId.ToUuid()
	if d.SongDao.Id == uuid.Nil {
		d.SongDao.Id = d.SongId
	}
	d.SongDao.FromDomain(value.SongData())
	d.SongId = d.SongDao.Id
	d.ScheduledAtUtc = value.ScheduledAtUtc().UTC()
}
