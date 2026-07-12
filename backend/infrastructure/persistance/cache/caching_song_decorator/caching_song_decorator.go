package caching_song_decorator

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache"
)

type CachingSongDecorator struct {
	decorated  i_music_data_provider_service.IMusicDataProviderService
	cache      cache.ICache[*music_data_response_dto.MusicDataResponseDto]
	unitOfWork i_unit_of_work.IUnitOfWork
}

func NewCachingSongDecorator(
	decorated i_music_data_provider_service.IMusicDataProviderService,
	cache cache.ICache[*music_data_response_dto.MusicDataResponseDto],
	unitOfWork i_unit_of_work.IUnitOfWork,
) *CachingSongDecorator {
	return &CachingSongDecorator{
		decorated:  decorated,
		cache:      cache,
		unitOfWork: unitOfWork,
	}
}

func (c *CachingSongDecorator) SearchSongsByQuery(ctx context.Context, accessToken, query string, nextPageToken *string, pageSize uint8) (*music_data_response_dto.MusicDataResponseDto, error) {
	result, cacheErr := c.cache.Get(query, ctx, c.unitOfWork.GetQueryer())
	if cacheErr == nil {
		return result, nil
	}

	result, err := c.decorated.SearchSongsByQuery(ctx, accessToken, query, nextPageToken, pageSize)
	if err != nil {
		return nil, err
	}

	cacheErr = c.cache.Set(query, result, ctx, c.unitOfWork.GetQueryer())
	if cacheErr != nil {
		return nil, cacheErr
	}

	return result, nil
}
