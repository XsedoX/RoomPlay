package default_playlist

import (
	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
)

type DefaultPlaylist struct {
	externalId  string
	userId      user_id.UserId
	songsAmount uint
	title       string
}

func (dp *DefaultPlaylist) ExternalId() string {
	return dp.externalId
}

func (dp *DefaultPlaylist) UserId() user_id.UserId {
	return dp.userId
}

func (dp *DefaultPlaylist) SongsAmount() uint {
	return dp.songsAmount
}

func (dp *DefaultPlaylist) Title() string {
	return dp.title
}

func NewDefaultPlaylist(
	externalId string,
	userId user_id.UserId,
	songsAmount uint,
	title string,
) (*DefaultPlaylist, error) {
	if externalId == "" {
		return nil, domain_errors.NewDefaultPlaylistExternalIdEmptyError()
	}
	if songsAmount == 0 {
		return nil, domain_errors.NewDefaultPlaylistSongsAmountZeroError()
	}
	if title == "" {
		return nil, domain_errors.NewDefaultPlaylistTitleEmptyError()
	}
	return &DefaultPlaylist{
		title:       title,
		externalId:  externalId,
		userId:      userId,
		songsAmount: songsAmount,
	}, nil
}

func (dp DefaultPlaylist) Equal(o DefaultPlaylist) bool {
	return dp.externalId == o.externalId &&
		dp.userId == o.userId &&
		dp.songsAmount == o.songsAmount &&
		dp.title == o.title
}
