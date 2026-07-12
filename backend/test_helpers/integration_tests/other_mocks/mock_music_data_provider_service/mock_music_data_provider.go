package mock_music_data_provider_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/stretchr/testify/mock"
)

type MockMusicDataProviderService struct {
	mock.Mock
}

func (m *MockMusicDataProviderService) SearchSongsByQuery(
	ctx context.Context,
	accessToken,
	query string,
	nextPageToken *string,
	pageSize uint8,
) (
	*music_data_response_dto.MusicDataResponseDto,
	error,
) {
	args := m.Called(ctx, accessToken, query, nextPageToken, pageSize)
	return args.Get(0).(*music_data_response_dto.MusicDataResponseDto), args.Error(1)
}
