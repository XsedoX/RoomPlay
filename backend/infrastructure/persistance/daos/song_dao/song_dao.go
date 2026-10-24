package song_dao

import (
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/song_external_data_dao"
	"github.com/google/uuid"
)

type SongDao struct {
	Id           uuid.UUID                                    `gorm:"primaryKey;column:id;type:uuid"`
	Title        string                                       `gorm:"column:title"`
	Author       string                                       `gorm:"column:author"`
	Isrc         *string                                      `gorm:"column:isrc"`
	ExternalData []song_external_data_dao.SongExternalDataDao `gorm:"foreignKey:SongId;references:Id"`
}

func (SongDao) TableName() string {
	return "songs"
}

// ToDomain requires ExternalData for the requested provider to be loaded.
// SongData describes one provider, whereas a SongDao may contain several.
func (d SongDao) ToDomain(provider music_provider.MusicProvider) *song_data.SongData {
	for _, externalData := range d.ExternalData {
		if externalData.MusicProvider == provider.String() {
			return externalData.ToDomain(d.Title, d.Author, d.Isrc)
		}
	}
	panic("SongDao.ToDomain: missing external data for " + provider.String())
}

// FromDomain updates this song's selected provider without discarding other providers.
// Populate Id (and any existing ExternalData) before calling when updating a stored
// song. A zero Id denotes a new row; resolving an existing song is the repository's job.
func (d *SongDao) FromDomain(value song_data.SongData) {
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	d.Title = value.Title()
	d.Author = value.Author()
	d.Isrc = nil
	if value.Isrc() != nil {
		d.Isrc = new(*value.Isrc())
	}

	for i := range d.ExternalData {
		if d.ExternalData[i].MusicProvider == value.MusicProvider().String() {
			d.ExternalData[i].FromDomain(value, d.Id)
			return
		}
	}
	var externalData song_external_data_dao.SongExternalDataDao
	externalData.FromDomain(value, d.Id)
	d.ExternalData = append(d.ExternalData, externalData)
}
