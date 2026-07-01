package search_song_query_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	search_song_query_response "github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query_query_response"
)

type SearchSongQueryHandler struct {
	unitOfWork                    i_unit_of_work.IUnitOfWork
	musicService                  i_music_data_provider_service.IMusicDataProviderService
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
}

func NewSearchSongQueryHandler(unitOfWork i_unit_of_work.IUnitOfWork,
	musicService i_music_data_provider_service.IMusicDataProviderService,
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository,
) *SearchSongQueryHandler {
	return &SearchSongQueryHandler{
		unitOfWork:                    unitOfWork,
		musicService:                  musicService,
		externalCredentialsRepository: externalCredentialsRepository,
	}
}

func (handler *SearchSongQueryHandler) Handle(ctx context.Context, query *search_song_query.SearchSongQuery) (*[]search_song_query_response.SearchSongQueryResponse, error) {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application_helpers.NewMissingUserIdInContextError
	}
	var response []search_song_query_response.SearchSongQueryResponse
	err := handler.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		accessToken, accessTokenError := handler.externalCredentialsRepository.AccessTokenByUserId(ctx, *userId, handler.unitOfWork.GetQueryer())
		if accessTokenError != nil {
			return custom_error.NewCustomError("SearchSongQueryHandler.GetAccessTokenByUserId",
				"Problem with getting access token for music service.",
				accessTokenError,
				custom_error_type.Unexpected)
		}
		songs, searchSongsErr := handler.musicService.SearchSongsByQuery(
			ctx,
			accessToken,
			query.Query,
			query.NextPageToken,
			query.PageSize,
		)
		if searchSongsErr != nil {
			return custom_error.NewCustomError("SearchSongQueryHandler.SearchSongsByQuery",
				"Problem with searching songs from music service.",
				searchSongsErr,
				custom_error_type.Unexpected)
		}
		for _, song := range *songs {
			response = append(response, search_song_query_response.SearchSongQueryResponse{
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
