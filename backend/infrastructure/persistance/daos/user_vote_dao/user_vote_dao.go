package user_vote_dao

import (
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_id"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/google/uuid"
)

type UserVoteDao struct {
	UserId         uuid.UUID `gorm:"primaryKey;foreignKey:UserId;references:Id"`
	EnqueuedSongId uuid.UUID `gorm:"primaryKey;foreignKey:EnqueuedSongId;references:Id"`
	State          string    `gorm:"column:vote_status;type:vote_status;default:'not_voted'" db:"vote_status"`
}

func (UserVoteDao) TableName() string {
	return "users_votes"
}

func (d *UserVoteDao) FromDomain(value vote_status.VoteStatus, userId user_id.UserId, enqueuedSongId enqueued_song_id.EnqueuedSongId) {
	*d = UserVoteDao{
		UserId:         userId.ToUuid(),
		EnqueuedSongId: *enqueuedSongId.ToUuid(),
		State:          value.String(),
	}
}
