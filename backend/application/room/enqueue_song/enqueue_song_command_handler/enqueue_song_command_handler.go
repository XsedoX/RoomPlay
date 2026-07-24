package enqueue_song_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	enquque_song_command "github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
)

type EnqueueSongCommandHandler struct {
	unitOfWork                    i_unit_of_work.IUnitOfWork
	roomRepository                i_room_repository.IRoomRepository
	musicDataProviderService      i_music_data_provider_service.IMusicDataProviderService
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
}

func NewEnqueueSongCommandHandler(
	unitOfWork i_unit_of_work.IUnitOfWork,
	roomRepository i_room_repository.IRoomRepository,
	musicDataProviderService i_music_data_provider_service.IMusicDataProviderService,
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository,
) *EnqueueSongCommandHandler {
	return &EnqueueSongCommandHandler{
		unitOfWork:                    unitOfWork,
		externalCredentialsRepository: externalCredentialsRepository,
		musicDataProviderService:      musicDataProviderService,
		roomRepository:                roomRepository,
	}
}

func (e *EnqueueSongCommandHandler) Handle(ctx context.Context, command *enquque_song_command.EnqueueSongCommand) error {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return application_helpers.NewMissingUserIdInContextError
	}

	err := e.unitOfWork.ExecuteTransaction(ctx, func(ctx context.Context) error {
		roomInstance, err := e.roomRepository.GetRoomAggregareByUserId(ctx, *userId, e.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		accessToken, accessTokenErr := e.externalCredentialsRepository.AccessTokenByUserId(ctx, *userId, e.unitOfWork.GetQueryer())
		if accessTokenErr != nil {
			return accessTokenErr
		}
		dataExternalSong, songErr := e.musicDataProviderService.GetSongById(ctx, accessToken, command.SongId)
		if songErr != nil {
			return songErr
		}
		songData, songDataErr := song_data.NewSongData(
			command.SongId,
			dataExternalSong.Title,
			dataExternalSong.Author,
			dataExternalSong.AlbumCoverUrl,
			dataExternalSong.LengthSeconds,
			dataExternalSong.MusicProvider,
			dataExternalSong.Isrc,
		)
		if songDataErr != nil {
			return songDataErr
		}
		roomInstance.EnqueueSong(*userId, *songData)
		err = e.roomRepository.UpdateRoom(ctx, roomInstance, e.unitOfWork.GetQueryer())
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
