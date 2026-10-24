package room_dao

import (
	"slices"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_encrypter"
	"github.com/XsedoX/RoomPlay/domain/external_credentials/music_provider"
	"github.com/XsedoX/RoomPlay/domain/room"
	"github.com/XsedoX/RoomPlay/domain/room/default_playlist"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song"
	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/room/scheduled_song"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/banned_user_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/default_playlist_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/enqueued_song_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/scheduled_song_dao"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos/user_room_data_dao"
	"github.com/google/uuid"
)

type RoomDao struct {
	Id                   uuid.UUID                                `gorm:"primaryKey"`
	Name                 string                                   `gorm:"column:name"`
	Password             []byte                                   `gorm:"column:password;type:bytea"`
	QrCodeHash           []byte                                   `gorm:"column:qr_code_hash;type:bytea"`
	BoostCooldownSeconds *uint16                                  `gorm:"column:boost_cooldown_seconds;type:smallint"`
	CreatedAtUtc         time.Time                                `gorm:"column:created_at_utc;type:timestamp"`
	LifespanSeconds      uint32                                   `gorm:"column:lifespan_seconds;type:int"`
	EnqueuedSongs        []enqueued_song_dao.EnqueuedSongDao      `gorm:"foreignKey:RoomId;references:Id"`
	UsersRoomData        []user_room_data_dao.UserRoomDataDao     `gorm:"foreignKey:RoomId;references:Id"`
	BannedUsers          []banned_user_dao.BannedUserDao          `gorm:"foreignKey:RoomId;references:Id"`
	ScheduledSong        *scheduled_song_dao.ScheduledSongDao     `gorm:"foreignKey:RoomId;references:Id"`
	DefaultPlaylist      *default_playlist_dao.DefaultPlaylistDao `gorm:"column:roomId;foreignKey:RoomId;references:Id"`
}

func (RoomDao) TableName() string {
	return "rooms"
}

// ToDomain expects the room's associations, including song provider data and
// individual votes, to be loaded. Only QR-code decryption can return an error.
func (d RoomDao) ToDomain(provider music_provider.MusicProvider, encrypter i_encrypter.IEncrypter) (*room.Room, error) {
	qrCode, err := encrypter.Decrypt(d.QrCodeHash)
	if err != nil {
		return nil, err
	}

	songs := make([]enqueued_song.EnqueuedSong, len(d.EnqueuedSongs))
	for i := range d.EnqueuedSongs {
		songs[i] = *d.EnqueuedSongs[i].ToDomain(provider)
	}
	members := make([]user_id.UserId, len(d.UsersRoomData))
	for i := range d.UsersRoomData {
		members[i] = user_id.UserId(d.UsersRoomData[i].UserId)
	}
	bannedUsers := make([]user_id.UserId, len(d.BannedUsers))
	for i := range d.BannedUsers {
		bannedUsers[i] = user_id.UserId(d.BannedUsers[i].UserId)
	}

	var boostCooldownSeconds *uint16
	if d.BoostCooldownSeconds != nil {
		boostCooldownSeconds = new(*d.BoostCooldownSeconds)
	}
	var scheduledSong *scheduled_song.ScheduledSong
	if d.ScheduledSong != nil {
		scheduledSong = d.ScheduledSong.ToDomain(provider)
	}
	var defaultPlaylist *default_playlist.DefaultPlaylist
	if d.DefaultPlaylist != nil {
		defaultPlaylist = d.DefaultPlaylist.ToDomain()
	}

	return room.HydrateRoom(
		room_id.RoomId(d.Id),
		d.Name,
		slices.Clone(d.Password),
		qrCode,
		boostCooldownSeconds,
		d.CreatedAtUtc,
		d.LifespanSeconds,
		songs,
		members,
		bannedUsers,
		scheduledSong,
		defaultPlaylist,
	), nil
}

// FromDomain needs the host ID because Room only stores member IDs, not roles.
// Mapping back into a loaded DAO retains matching queue entries' song IDs/vote
// rows and members' boost timestamps. A fresh DAO cannot reconstruct that metadata.
// Existing shared songs must be resolved by the repository before persistence.
func (d *RoomDao) FromDomain(value *room.Room, hostId user_id.UserId, encrypter i_encrypter.IEncrypter) error {
	qrCodeHash, err := encrypter.Encrypt(value.QrCode())
	if err != nil {
		return err
	}

	result := RoomDao{
		Id:              value.Id().ToUuid(),
		Name:            value.Name(),
		Password:        slices.Clone(value.Password()),
		QrCodeHash:      qrCodeHash,
		CreatedAtUtc:    value.CreatedAtUtc(),
		LifespanSeconds: value.LifespanSeconds(),
	}
	if value.BoostCooldownSeconds() != nil {
		result.BoostCooldownSeconds = new(*value.BoostCooldownSeconds())
	}

	existingSongs := make(map[uuid.UUID]enqueued_song_dao.EnqueuedSongDao, len(d.EnqueuedSongs))
	for _, existingSong := range d.EnqueuedSongs {
		existingSongs[existingSong.Id] = existingSong
	}
	songs := value.AllSongs()
	result.EnqueuedSongs = make([]enqueued_song_dao.EnqueuedSongDao, len(songs))
	for i := range songs {
		enqueuedSongId := songs[i].Id()
		song := existingSongs[*enqueuedSongId.ToUuid()]
		song.FromDomain(&songs[i], value.Id())
		result.EnqueuedSongs[i] = song
	}

	existingMemberships := make(map[uuid.UUID]user_room_data_dao.UserRoomDataDao, len(d.UsersRoomData))
	for _, membership := range d.UsersRoomData {
		if membership.RoomId == result.Id {
			existingMemberships[membership.UserId] = membership
		}
	}
	members := value.Members()
	result.UsersRoomData = make([]user_room_data_dao.UserRoomDataDao, len(members))
	for i, memberId := range members {
		role := user_role.Member
		if memberId == hostId {
			role = user_role.Host
		}
		membership := existingMemberships[memberId.ToUuid()]
		result.UsersRoomData[i].FromDomain(memberId, value.Id(), role, membership.BoostUsedAtUtc)
	}

	bannedUsers := value.BannedUsers()
	result.BannedUsers = make([]banned_user_dao.BannedUserDao, len(bannedUsers))
	for i, bannedUserId := range bannedUsers {
		result.BannedUsers[i].FromDomain(bannedUserId, value.Id())
	}

	if value.ScheduledSong() != nil {
		scheduledSong := scheduled_song_dao.ScheduledSongDao{}
		if d.ScheduledSong != nil {
			// Reuse a song ID only when the scheduled track is still the same.
			songData := value.ScheduledSong().SongData()
			for _, externalData := range d.ScheduledSong.SongDao.ExternalData {
				if externalData.MusicProvider == songData.MusicProvider().String() && externalData.ExternalId == songData.ExternalId() {
					scheduledSong = *d.ScheduledSong
					break
				}
			}
		}
		scheduledSong.FromDomain(value.ScheduledSong(), value.Id())
		result.ScheduledSong = &scheduledSong
	}
	if value.DefaultPlaylist() != nil {
		result.DefaultPlaylist = new(default_playlist_dao.DefaultPlaylistDao)
		result.DefaultPlaylist.FromDomain(value.DefaultPlaylist(), value.Id())
	}

	*d = result
	return nil
}
