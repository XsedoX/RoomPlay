package song_external_data_dao

import (
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/google/uuid"
)

type SongExternalDataDao struct {
	Id            uuid.UUID `gorm:"primaryKey;column:id" db:"id"`
	SongId        uuid.UUID `gorm:"column:song_id;foreignKey:SongId;references:Id" db:"song_id"`
	LengthSeconds uint16    `gorm:"column:length_seconds;type:smallint" db:"length_seconds"`
	AlbumCoverUrl string    `gorm:"column:album_cover_url;type:text" db:"album_cover_url"`
	ExternalId    string    `gorm:"column:external_id;type:text" db:"external_id"`
	MusicProvider string    `gorm:"column:music_provider;type:music_provider" db:"music_provider"`
}

func (SongExternalDataDao) TableName() string {
	return "songs_external_data"
}

// ToDomain combines provider data with the fields stored on SongDao.
func (d SongExternalDataDao) ToDomain(title, author string, isrc *string) *song_data.SongData {
	var domainIsrc *string
	if isrc != nil {
		domainIsrc = new(*isrc)
	}
	return song_data.HydrateSongData(
		d.ExternalId,
		title,
		author,
		d.AlbumCoverUrl,
		d.LengthSeconds,
		*music_provider.ParseMusicProvider(d.MusicProvider),
		domainIsrc,
	)
}

// FromDomain retains an existing row ID; SongData has no persistence ID of its own.
func (d *SongExternalDataDao) FromDomain(value song_data.SongData, songId uuid.UUID) {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	d.SongId = songId
	d.LengthSeconds = value.LengthSeconds()
	d.AlbumCoverUrl = value.AlbumCoverUrl()
	d.ExternalId = value.ExternalId()
	d.MusicProvider = value.MusicProvider().String()
}
