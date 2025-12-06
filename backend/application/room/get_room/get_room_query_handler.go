package get_room

import (
	"context"
	"encoding/base64"

	"xsedox.com/main/application"
	"xsedox.com/main/application/contracts"
	"xsedox.com/main/application/custom_errors"
	contracts3 "xsedox.com/main/application/room/contracts"
)

type GetRoomQueryHandler struct {
	unitOfWork     contracts.IUnitOfWork
	roomRepository contracts3.IRoomRepository
}

func NewGetRoomQueryHandler(unitOfWork contracts.IUnitOfWork,
	roomRepository contracts3.IRoomRepository) *GetRoomQueryHandler {
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
			return custom_errors.NewCustomError("GetRoomQueryHandler.GetRoomByUserId",
				"Couldn't get user's room.",
				getRoomDataErr,
				custom_errors.Unexpected)
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

		if roomData.PlayingSongTitle != nil && roomData.PlayingSongAuthor != nil {
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
