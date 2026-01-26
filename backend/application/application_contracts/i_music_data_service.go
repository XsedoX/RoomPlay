package application_contracts

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/dtos"
)

type IMusicDataService interface {
	SearchSongsByQuery(ctx context.Context, accessToken, query string) (*[]dtos.MusicDataResponseDto, error)
}
