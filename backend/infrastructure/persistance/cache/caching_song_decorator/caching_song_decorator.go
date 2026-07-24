package caching_song_decorator

import (
	"context"
	"log"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/cache"
)

type CachingSongDecorator struct {
	decorated             i_music_data_provider_service.IMusicDataProviderService
	cache                 cache.ICache[*music_data_response_dto.MusicDataResponseDto]
	unitOfWork            i_unit_of_work.IUnitOfWork
	songByExternalIdCache cache.ICache[*music_data_response_dto.SongDataResponseDto]
}

func NewCachingSongDecorator(
	decorated i_music_data_provider_service.IMusicDataProviderService,
	cache cache.ICache[*music_data_response_dto.MusicDataResponseDto],
	unitOfWork i_unit_of_work.IUnitOfWork,
	songByExternalIdCache cache.ICache[*music_data_response_dto.SongDataResponseDto],
) *CachingSongDecorator {
	return &CachingSongDecorator{
		decorated:             decorated,
		cache:                 cache,
		unitOfWork:            unitOfWork,
		songByExternalIdCache: songByExternalIdCache,
	}
}

func (c *CachingSongDecorator) GetSongById(ctx context.Context, accessToken, songId string) (*music_data_response_dto.SongDataResponseDto, error) {
	result := &music_data_response_dto.SongDataResponseDto{}
	cacheErr := c.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		var err error
		result, err = c.songByExternalIdCache.GetExact(songId, ctx, c.unitOfWork.GetQueryer())
		return err
	})
	if cacheErr == nil {
		return result, nil
	}

	result, err := c.decorated.GetSongById(ctx, accessToken, songId)
	if err != nil {
		return nil, err
	}

	cacheErr = c.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		cacheErr := c.songByExternalIdCache.Set(
			songId,
			result,
			ctx,
			c.unitOfWork.GetQueryer(),
		)
		return cacheErr
	})
	if cacheErr != nil {
		return nil, cacheErr
	}

	return result, nil
}

func (c *CachingSongDecorator) SearchSongsByQuery(ctx context.Context, accessToken, query string, nextPageToken *string, pageSize uint8) (*music_data_response_dto.MusicDataResponseDto, error) {
	cacheKey := query
	if nextPageToken != nil {
		cacheKey += "::" + *nextPageToken
	}

	var result *music_data_response_dto.MusicDataResponseDto
	cacheErr := c.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		var err error
		result, err = c.cache.Get(cacheKey, ctx, c.unitOfWork.GetQueryer())
		return err
	})
	if cacheErr == nil {
		return result, nil
	}

	result, err := c.decorated.SearchSongsByQuery(ctx, accessToken, query, nextPageToken, pageSize)
	if err != nil {
		return nil, err
	}
	_ = c.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		for _, song := range result.Songs {
			songInstance := &song
			err := c.songByExternalIdCache.Set(
				song.VideoId,
				songInstance,
				ctx,
				c.unitOfWork.GetQueryer(),
			)
			if err != nil {
				log.Printf("Error caching song by external ID: %v", err)
			}
		}
		return nil
	})

	cacheErr = c.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		cacheErr := c.cache.Set(cacheKey, result, ctx, c.unitOfWork.GetQueryer())
		return cacheErr
	})
	if cacheErr != nil {
		return nil, cacheErr
	}

	return result, nil
}
