package get_room_query_handler

import (
	"context"
	"encoding/base64"

	"github.com/XsedoX/RoomPlay/application/application_contracts/i_unit_of_work"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/custom_error/custom_error_type"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts/i_room_repository"
)

type GetRoomQueryHandler struct {
	unitOfWork     i_unit_of_work.IUnitOfWork
	roomRepository i_room_repository.IRoomRepository
}

func NewGetRoomQueryHandler(unitOfWork i_unit_of_work.IUnitOfWork,
	roomRepository i_room_repository.IRoomRepository,
) *GetRoomQueryHandler {
	return &GetRoomQueryHandler{
		unitOfWork:     unitOfWork,
		roomRepository: roomRepository,
	}
}

func (r GetRoomQueryHandler) Handle(ctx context.Context) (*get_room_query_response.GetRoomQueryResponse, error) {
	var response get_room_query_response.GetRoomQueryResponse
	userId, ok := application_helpers.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application_helpers.NewMissingUserIdInContextError
	}
	err := r.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		roomData, getRoomDataErr := r.roomRepository.GetRoomByUserId(ctx, *userId, r.unitOfWork.GetQueryer())
		if getRoomDataErr != nil {
			return custom_error.NewCustomError("GetRoomQueryHandler.GetRoomByUserId",
				"Couldn't get user's room.",
				getRoomDataErr,
				custom_error_type.Unexpected)
		}
		response.Name = roomData.Name
		response.QrCode = base64.RawURLEncoding.EncodeToString(roomData.QrCode)
		response.UserRole = roomData.UserRole
		if roomData.BoostUsedAtUtc != nil && roomData.BoostCooldownSeconds != nil {
			response.BoostData = &get_room_query_response.BoostDataDto{
				BoostUsedAtUtc:       *roomData.BoostUsedAtUtc,
				BoostCooldownSeconds: *roomData.BoostCooldownSeconds,
			}
		} else {
			response.BoostData = nil
		}

		if roomData.PlayingSongTitle != nil && roomData.PlayingSongAuthor != nil && roomData.PlayingSongStartedAtUtc != nil && roomData.PlayingSongLengthSeconds != nil {
			response.PlayingSong = &get_room_query_response.PlayingSongDto{
				Title:         *roomData.PlayingSongTitle,
				Author:        *roomData.PlayingSongAuthor,
				StartedAtUtc:  *roomData.PlayingSongStartedAtUtc,
				LengthSeconds: *roomData.PlayingSongLengthSeconds,
			}
		} else {
			response.PlayingSong = nil
		}
		response.Songs = make([]get_room_query_response.RoomSongListDto, 0)
		for _, songDb := range roomData.SongDaos {
			response.Songs = append(response.Songs, get_room_query_response.RoomSongListDto{
				Title:         songDb.Title,
				Author:        songDb.Author,
				AddedBy:       songDb.AddedBy,
				Votes:         songDb.Votes,
				AlbumCoverUrl: songDb.AlbumCoverUrl,
				Id:            songDb.Id,
				State:         songDb.State,
				VoteStatus:    songDb.VoteStatus,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}
