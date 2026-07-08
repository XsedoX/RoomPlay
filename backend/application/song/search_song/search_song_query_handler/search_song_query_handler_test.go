package search_song_query_handler

import (
	"context"
	"testing"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
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
		musicProviderResponse := []music_data_response_dto.MusicDataResponseDto{
			{
				VideoId:       gofakeit.ID(),
				Title:         gofakeit.SongName(),
				Author:        gofakeit.SongArtist(),
				AlbumCoverUrl: gofakeit.URL(),
				NextPageToken: gofakeit.ID(),
			},
			{
				VideoId:       gofakeit.ID(),
				Title:         gofakeit.SongName(),
				Author:        gofakeit.SongArtist(),
				AlbumCoverUrl: gofakeit.URL(),
				NextPageToken: gofakeit.ID(),
			},
			{
				VideoId:       gofakeit.ID(),
				Title:         gofakeit.SongName(),
				Author:        gofakeit.SongArtist(),
				AlbumCoverUrl: gofakeit.URL(),
				NextPageToken: gofakeit.ID(),
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

		require.Equal(t, response[0].VideoId, musicProviderResponse[0].VideoId)
		require.Equal(t, response[1].VideoId, musicProviderResponse[1].VideoId)
		require.Equal(t, response[2].VideoId, musicProviderResponse[2].VideoId)

		require.Equal(t, response[0].Title, musicProviderResponse[0].Title)
		require.Equal(t, response[1].Title, musicProviderResponse[1].Title)
		require.Equal(t, response[2].Title, musicProviderResponse[2].Title)

		require.Equal(t, response[0].Author, musicProviderResponse[0].Author)
		require.Equal(t, response[1].Author, musicProviderResponse[1].Author)
		require.Equal(t, response[2].Author, musicProviderResponse[2].Author)

		require.Equal(t, response[0].AlbumCoverUrl, musicProviderResponse[0].AlbumCoverUrl)
		require.Equal(t, response[1].AlbumCoverUrl, musicProviderResponse[1].AlbumCoverUrl)
		require.Equal(t, response[2].AlbumCoverUrl, musicProviderResponse[2].AlbumCoverUrl)

		require.Equal(t, response[0].NextPageToken, musicProviderResponse[0].NextPageToken)
		require.Equal(t, response[1].NextPageToken, musicProviderResponse[1].NextPageToken)
		require.Equal(t, response[2].NextPageToken, musicProviderResponse[2].NextPageToken)
	})
}
