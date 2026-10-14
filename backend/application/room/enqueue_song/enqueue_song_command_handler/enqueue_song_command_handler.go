package enqueue_song_command_handler

import (
	"context"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_event_publisher"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_authentication_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_external_credentials_repository"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	enquque_song_command "github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/song_data"
	"github.com/XsedoX/RoomPlay/domain/shared"
)

type EnqueueSongCommandHandler struct {
	unitOfWork                    i_unit_of_work.IUnitOfWork
	roomRepository                i_room_repository.IRoomRepository
	musicDataProviderService      i_music_data_provider_service.IMusicDataProviderService
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository
	authenticationService         i_external_authentication_service.IExternalAuthenticationService
	domainEventsPublisher         i_event_publisher.IEventPublisher
}

func NewEnqueueSongCommandHandler(
	unitOfWork i_unit_of_work.IUnitOfWork,
	roomRepository i_room_repository.IRoomRepository,
	musicDataProviderService i_music_data_provider_service.IMusicDataProviderService,
	externalCredentialsRepository i_external_credentials_repository.IExternalCredentialsRepository,
	authenticationService i_external_authentication_service.IExternalAuthenticationService,
	domainEventsPublisher i_event_publisher.IEventPublisher,
) *EnqueueSongCommandHandler {
	return &EnqueueSongCommandHandler{
		unitOfWork:                    unitOfWork,
		authenticationService:         authenticationService,
		externalCredentialsRepository: externalCredentialsRepository,
		musicDataProviderService:      musicDataProviderService,
		roomRepository:                roomRepository,
		domainEventsPublisher:         domainEventsPublisher,
	}
}

func (e *EnqueueSongCommandHandler) Handle(ctx context.Context, command *enquque_song_command.EnqueueSongCommand) error {
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return application_helpers.NewMissingUserIdInContextError("EnqueueSongCommandHandler.Handle")
	}

	events := make([]shared.IDomainEvent, 0)

	err := e.unitOfWork.ExecuteTransaction(
		ctx,
		func(ctx context.Context) error {
			roomInstance, err := e.roomRepository.GetRoomAggregateByUserId(ctx, *userId, e.unitOfWork.GetQueryer(ctx))
			if err != nil {
				return err
			}
			accessToken, accessTokenErr := e.externalCredentialsRepository.GetAccessTokenByUserId(ctx, *userId, e.unitOfWork.GetQueryer(ctx))
			if accessTokenErr != nil {
				return accessTokenErr
			}

			var authErr error
			if accessToken.IsExpired() {
				accessToken, authErr = e.authenticationService.RefreshAccessTokenWithExternalProvider(ctx, *userId)
				if authErr != nil {
					return authErr
				}
			}
			dataExternalSong, songErr := e.musicDataProviderService.GetSongById(
				ctx,
				accessToken.Value(),
				command.SongExternalId,
			)
			if songErr != nil {
				return songErr
			}
			songData, songDataErr := song_data.NewSongData(
				command.SongExternalId,
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
			events = roomInstance.EnqueueSong(*userId, *songData)
			err = e.roomRepository.UpdateRoom(ctx, roomInstance, e.unitOfWork.GetQueryer(ctx))
			if err != nil {
				return err
			}
			return nil
		},
	)
	for _, event := range events {
		e.domainEventsPublisher.Publish(event)
	}

	return err
}
