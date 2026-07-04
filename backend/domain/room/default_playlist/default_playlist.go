package default_playlist

import (
	"github.com/XsedoX/RoomPlay/domain/domain_errors/empty_string_domain_error"
	"github.com/XsedoX/RoomPlay/domain/domain_errors/validation_domain_error"
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
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"DefaultPlaylist.ExternalId",
			"external id",
		)
	}
	if songsAmount == 0 {
		return nil, validation_domain_error.NewValidationDomainError(
			"DefaultPlaylist.SongsAmount.Zero",
			"Songs amount cannot be zero",
		)
	}
	if title == "" {
		return nil, empty_string_domain_error.NewEmptyStringDomainError(
			"DefaultPlaylist.Title",
			"title",
		)
	}
	return &DefaultPlaylist{
		title:       title,
		externalId:  externalId,
		userId:      userId,
		songsAmount: songsAmount,
	}, nil
}
