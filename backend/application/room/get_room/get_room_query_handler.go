package get_room

import (
	"context"
	"encoding/base64"

	"github.com/XsedoX/RoomPlay/application"
	"github.com/XsedoX/RoomPlay/application/application_contracts"
	"github.com/XsedoX/RoomPlay/application/customerrors"
	"github.com/XsedoX/RoomPlay/application/room/room_contracts"
)

type GetRoomQueryHandler struct {
	unitOfWork     application_contracts.IUnitOfWork
	roomRepository room_contracts.IRoomRepository
}

func NewGetRoomQueryHandler(unitOfWork application_contracts.IUnitOfWork,
	roomRepository room_contracts.IRoomRepository,
) *GetRoomQueryHandler {
	return &GetRoomQueryHandler{
		unitOfWork:     unitOfWork,
		roomRepository: roomRepository,
	}
}

func (r GetRoomQueryHandler) Handle(ctx context.Context) (*GetRoomQueryResponse, error) {
	var response GetRoomQueryResponse
	userId, ok := application.GetUserIdFromContext(ctx)
	if !ok {
		return nil, application.NewMissingUserIdInContextError
	}
	err := r.unitOfWork.ExecuteRead(ctx, func(ctx context.Context) error {
		roomData, getRoomDataErr := r.roomRepository.GetRoomByUserId(ctx, *userId, r.unitOfWork.GetQueryer())
		if getRoomDataErr != nil {
			return customerrors.NewCustomError("GetRoomQueryHandler.GetRoomByUserId",
				"Couldn't get user's room.",
				getRoomDataErr,
				customerrors.Unexpected)
		}
		response.Name = roomData.Name
		response.QrCode = base64.RawURLEncoding.EncodeToString(roomData.QrCode)
		response.UserRole = roomData.UserRole
		if roomData.BoostUsedAtUtc != nil && roomData.BoostCooldownSeconds != nil {
			response.BoostData = &BoostDataDto{
				BoostUsedAtUtc:       *roomData.BoostUsedAtUtc,
				BoostCooldownSeconds: *roomData.BoostCooldownSeconds,
			}
		} else {
			response.BoostData = nil
		}

		if roomData.PlayingSongTitle != nil && roomData.PlayingSongAuthor != nil && roomData.PlayingSongStartedAtUtc != nil && roomData.PlayingSongLengthSeconds != nil {
			response.PlayingSong = &PlayingSongDto{
				Title:         *roomData.PlayingSongTitle,
				Author:        *roomData.PlayingSongAuthor,
				StartedAtUtc:  *roomData.PlayingSongStartedAtUtc,
				LengthSeconds: *roomData.PlayingSongLengthSeconds,
			}
		} else {
			response.PlayingSong = nil
		}
		response.Songs = make([]RoomSongListDto, 0)
		for _, songDb := range roomData.SongDaos {
			response.Songs = append(response.Songs, RoomSongListDto{
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
