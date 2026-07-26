package search_song_query_handler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/application/song/search_song/search_song_query"
	"github.com/XsedoX/RoomPlay/domain/token"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache/caching_song_decorator"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_external_authentication_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_cache"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupMocks(t *testing.T) (
	mockUoW *mock_unit_of_work.MockUnitOfWork,
	mockMusicDataProvider *mock_music_data_provider_service.MockMusicDataProviderService,
	mockExternalCredentialsRepository *mock_external_credentials_repository.MockExternalCredentialsRepository,
	userId user_id.UserId,
	ctx context.Context,
	mockSongsCache *mock_cache.MockCache[*music_data_response_dto.MusicDataResponseDto],
	mockExternalSongsCache *mock_cache.MockCache[*music_data_response_dto.SongDataResponseDto],
	mockExternalAuthenticationService *mock_external_authentication_service.MockExternalAuthenticationService,
) {
	mockUoW = new(mock_unit_of_work.MockUnitOfWork)
	mockMusicDataProvider = new(mock_music_data_provider_service.MockMusicDataProviderService)
	userId, ctx = test_helpers.AddUserIdToContext(context.Background())
	mockExternalCredentialsRepository = new(mock_external_credentials_repository.MockExternalCredentialsRepository)
	mockSongsCache = new(mock_cache.MockCache[*music_data_response_dto.MusicDataResponseDto])
	mockExternalSongsCache = new(mock_cache.MockCache[*music_data_response_dto.SongDataResponseDto])
	mockExternalAuthenticationService = new(mock_external_authentication_service.MockExternalAuthenticationService)

	defer func() {
		mockUoW.AssertExpectations(t)
		mockMusicDataProvider.AssertExpectations(t)
		mockExternalCredentialsRepository.AssertExpectations(t)
		mockSongsCache.AssertExpectations(t)
		mockExternalSongsCache.AssertExpectations(t)
		mockExternalAuthenticationService.AssertExpectations(t)
	}()
	return
}

func TestSearchSongQueryHandlerCacheClear(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		mockUoW,
			mockMusicDataProvider,
			mockExternalCredentialsRepository,
			userId,
			ctx,
			mockSongsCache,
			mockExternalSongsCache,
			mockAuthService := setupMocks(t)

		accessToken := token.HydrateToken("access_token", time.Now().Add(time.Hour*1))
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
			"GetAccessTokenByUserId",
			ctx,
			userId,
			mockUoW.GetQueryer(ctx),
		).Return(accessToken, nil)

		mockMusicDataProvider.On(
			"SearchSongsByQuery",
			ctx,
			accessToken.Value(),
			queryString,
			(*string)(nil),
			uint8(3)).Return(musicProviderResponse, nil)

		mockSongsCache.On(
			"Get",
			queryString,
			ctx,
			mockUoW.GetQueryer(ctx),
		).Return(nil, sql.ErrNoRows)

		mockSongsCache.On(
			"Set",
			mock.Anything,
			mock.Anything,
			ctx,
			mockUoW.GetQueryer(ctx),
		).Return(nil)

		mockExternalSongsCache.On(
			"Set",
			mock.Anything,
			mock.Anything,
			ctx,
			mockUoW.GetQueryer(ctx),
		).Return(nil)

		mockAuthService.On(
			"RefreshAccessTokenWithExternalProvider",
			ctx,
			userId,
		).Return(accessToken, nil)

		cachingSongDecorator := caching_song_decorator.NewCachingSongDecorator(
			mockMusicDataProvider,
			mockSongsCache,
			mockUoW,
			mockExternalSongsCache,
		)

		handler := NewSearchSongQueryHandler(
			mockUoW,
			cachingSongDecorator,
			mockExternalCredentialsRepository,
			mockAuthService,
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
