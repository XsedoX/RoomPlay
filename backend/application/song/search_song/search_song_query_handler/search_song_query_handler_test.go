package search_song_query_handler

import (
	"context"
	"testing"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func setupMocks(t *testing.T) (
	*mock_unit_of_work.MockUnitOfWork,
	*mock_music_data_provider_service.MockMusicDataProviderService,
	*mock_external_credentials_repository.MockExternalCredentialsRepository,
	user_id.UserId,
	context.Context,
) {
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	mockMusicDataProvider := new(mock_music_data_provider_service.MockMusicDataProviderService)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())
	mockExternalCredentialsRepository := new(mock_external_credentials_repository.MockExternalCredentialsRepository)

	defer func() {
		mockUoW.AssertExpectations(t)
		mockMusicDataProvider.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
	}()

	return mockUoW, mockMusicDataProvider, mockExternalCredentialsRepository, userId, ctx
}

func TestSearchSongQueryHandler(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockUoW,
			mockMusicDataProvider,
			mockExternalCredentialsRepository,
			userId,
			ctx := setupMocks(t)

		accessToken := "access_token"
		queryString := "test query"
		musicProviderResponse := &music_data_response_dto.MusicDataResponseDto{
			Songs: []music_data_response_dto.SongDataResponseDto{
				{
					VideoId:       gofakeit.ID(),
					Title:         gofakeit.SongName(),
					Author:        gofakeit.SongArtist(),
					AlbumCoverUrl: gofakeit.URL(),
				},
				{
					VideoId:       gofakeit.ID(),
					Title:         gofakeit.SongName(),
					Author:        gofakeit.SongArtist(),
					AlbumCoverUrl: gofakeit.URL(),
				},
				{
					VideoId:       gofakeit.ID(),
					Title:         gofakeit.SongName(),
					Author:        gofakeit.SongArtist(),
					AlbumCoverUrl: gofakeit.URL(),
				},
			},
			PageMetaDto: page_meta_dto.PageMetaDto{
				NextPageToken:     new(gofakeit.ID()),
				PreviousPageToken: new(gofakeit.ID()),
				HasNextPage:       true,
				PageSize:          uint8(gofakeit.Number(1, 10)),
			},
		}

		mockUoW.On("GetQueryer").Return(nil)

		mockExternalCredentialsRepository.On(
			"AccessTokenByUserId",
			ctx,
			userId,
			mockUoW.GetQueryer(),
		).Return(accessToken, nil)

		mockMusicDataProvider.On(
			"SearchSongsByQuery",
			ctx,
			accessToken,
			queryString,
			(*string)(nil),
			uint8(3)).Return(musicProviderResponse, nil)

		handler := NewSearchSongQueryHandler(
			mockUoW,
			mockMusicDataProvider,
			mockExternalCredentialsRepository,
		)
		query := search_song_query.SearchSongQuery{
			Query:         "test query",
			NextPageToken: nil,
			PageSize:      3,
		}

		response, err := handler.Handle(ctx, &query)
		require.NoError(t, err)

		require.Equal(t, response.Songs[0].VideoId, musicProviderResponse.Songs[0].VideoId)
		require.Equal(t, response.Songs[1].VideoId, musicProviderResponse.Songs[1].VideoId)
		require.Equal(t, response.Songs[2].VideoId, musicProviderResponse.Songs[2].VideoId)

		require.Equal(t, response.Songs[0].Title, musicProviderResponse.Songs[0].Title)
		require.Equal(t, response.Songs[1].Title, musicProviderResponse.Songs[1].Title)
		require.Equal(t, response.Songs[2].Title, musicProviderResponse.Songs[2].Title)

		require.Equal(t, response.Songs[0].Author, musicProviderResponse.Songs[0].Author)
		require.Equal(t, response.Songs[1].Author, musicProviderResponse.Songs[1].Author)
		require.Equal(t, response.Songs[2].Author, musicProviderResponse.Songs[2].Author)

		require.Equal(t, response.Songs[0].AlbumCoverUrl, musicProviderResponse.Songs[0].AlbumCoverUrl)
		require.Equal(t, response.Songs[1].AlbumCoverUrl, musicProviderResponse.Songs[1].AlbumCoverUrl)
		require.Equal(t, response.Songs[2].AlbumCoverUrl, musicProviderResponse.Songs[2].AlbumCoverUrl)

		require.Equal(t, response.PageMetaDto.NextPageToken, musicProviderResponse.PageMetaDto.NextPageToken)
		require.Equal(t, response.PageMetaDto.PageSize, musicProviderResponse.PageMetaDto.PageSize)
		require.Equal(t, response.PageMetaDto.PreviousPageToken, musicProviderResponse.PageMetaDto.PreviousPageToken)
		require.Equal(t, response.PageMetaDto.HasNextPage, musicProviderResponse.PageMetaDto.HasNextPage)
	})
}
