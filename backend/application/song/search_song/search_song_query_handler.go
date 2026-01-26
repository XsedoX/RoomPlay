package search_song

import (
	"context"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
)

type SearchSongQueryHandler struct {
	unitOfWork                    application_contracts.IUnitOfWork
	musicService                  application_contracts.IMusicDataService
	externalCredentialsRepository application_contracts.IExternalCredentialsRepository
}

func NewSearchSongQueryHandler(unitOfWork application_contracts.IUnitOfWork,
	musicService application_contracts.IMusicDataService,
	externalCredentialsRepository application_contracts.IExternalCredentialsRepository,
) *SearchSongQueryHandler {
	return &SearchSongQueryHandler{
		unitOfWork:                    unitOfWork,
		musicService:                  musicService,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (handler *SearchSongQueryHandler) Handle(ctx context.Context, query *SearchSongQuery) (*[]SearchSongQueryResponse, error) {
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application.NewMissingUserIdInContextError
	}
	var response []SearchSongQueryResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		accessToken, accessTokenError := handler.externalCredentialsRepository.GetAccessTokenByUserId(ctx, *userId, handler.unitOfWork.GetQueryer())
		if accessTokenError != nil {
			return customerrors.NewCustomError("SearchSongQueryHandler.GetAccessTokenByUserId",
				"Problem with getting access token for music service.",
				accessTokenError,
				customerrors.Unexpected)
		}
		songs, searchSongsErr := handler.musicService.SearchSongsByQuery(ctx, accessToken, query.Query)
		if searchSongsErr != nil {
			return customerrors.NewCustomError("SearchSongQueryHandler.SearchSongsByQuery",
				"Problem with searching songs from music service.",
				searchSongsErr,
				customerrors.Unexpected)
		}
		for _, song := range *songs {
			response = append(response, SearchSongQueryResponse{
				Url:           song.Url,
				Author:        song.Author,
				AlbumCoverUrl: song.AlabumCoverUrl,
				Title:         song.Title,
				LengthSeconds: song.LengthSeconds,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
