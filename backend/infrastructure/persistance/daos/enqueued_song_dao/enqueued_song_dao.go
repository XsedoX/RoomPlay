package enqueued_song_dao

import (
	"time"

	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/song_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_vote_dao"
	"github.com/google/uuid"
)

type EnqueuedSongDao struct {
	Id             uuid.UUID                   `gorm:"primaryKey"`
	RoomId         uuid.UUID                   `gorm:"column:room_id;"`
	SongId         uuid.UUID                   `gorm:"column:song_id;foreignKey:SongId;references:Id"`
	SongDao        song_dao.SongDao            `gorm:"foreignKey:SongId;references:Id"`
	AddedBy        uuid.UUID                   `gorm:"column:added_by;foreignKey:AddedBy;references:Id"`
	AddedByUserDao user_dao.UserDao            `gorm:"foreignKey:AddedBy;references:Id"`
	AddeedAtUtc    time.Time                   `gorm:"column:added_at_utc;type:timestamp"`
	StartedAtUtc   *time.Time                  `gorm:"column:started_at_utc;type:timestamp"`
	State          string                      `gorm:"column:state;type:song_state"`
	UserVotes      []user_vote_dao.UserVoteDao `gorm:"foreignKey:EnqueuedSongId;references:Id"`
}

func (EnqueuedSongDao) TableName() string {
	return "enqueued_songs"
}

// ToDomain expects SongDao.ExternalData and UserVotes to be loaded.
func (d EnqueuedSongDao) ToDomain(provider music_provider.MusicProvider) *enqueued_song.EnqueuedSong {
	votes := 0
	for _, vote := range d.UserVotes {
		switch *vote_status.ParseVoteStatus(vote.State) {
		case vote_status.Upvoted:
			votes++
		case vote_status.Downvoted:
			votes--
		}
	}
	// The domain stores an int8 score; do not silently wrap a larger total.
	if votes < -128 || votes > 127 {
		panic("EnqueuedSongDao.ToDomain: vote total exceeds the domain int8 range")
	}
	return enqueued_song.HydrateEnqueuedSong(
		enqueued_song_id.EnqueuedSongId(d.Id),
		*d.SongDao.ToDomain(provider),
		d.AddeedAtUtc,
		d.StartedAtUtc,
		*enqueued_song_state.ParseSongState(d.State),
		int8(votes),
		user_id.UserId(d.AddedBy),
	)
}

// FromDomain retains loaded vote rows for the same queue entry. A score alone
// cannot reconstruct who voted; new vote rows must be mapped separately.
// Existing song persistence IDs can be supplied through SongDao/SongId.
func (d *EnqueuedSongDao) FromDomain(value *enqueued_song.EnqueuedSong, roomId room_id.RoomId) {
	enqueuedSongId := value.Id()
	if d.Id != *enqueuedSongId.ToUuid() {
		d.UserVotes = nil
	}
	// This identifies the queue entry, not the underlying song.
	d.Id = *enqueuedSongId.ToUuid()
	d.RoomId = roomId.ToUuid()

	// The song's identity comes only from the song DAO or its foreign key.
	if d.SongDao.Id == uuid.Nil {
		d.SongDao.Id = d.SongId
	}
	d.SongDao.FromDomain(value.SongData())
	d.SongId = d.SongDao.Id

	d.AddedBy = value.AddedBy().ToUuid()
	d.AddedByUserDao = user_dao.UserDao{}
	d.AddeedAtUtc = value.AddedAtUtc()
	d.StartedAtUtc = value.StartedAtUtc()
	d.State = value.State().String()
}
