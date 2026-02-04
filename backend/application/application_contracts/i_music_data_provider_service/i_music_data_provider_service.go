package i_music_data_provider_service

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
)

type IMusicDataProviderService interface {
	SearchSongsByQuery(ctx context.Context, accessToken, query string) (*[]music_data_response_dto.MusicDataResponseDto, error)
}
