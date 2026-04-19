package get_room_query_handler

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/custom_error"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_song_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpMocks(t *testing.T) (*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*user_id.UserId,
	context.Context,
) {
	mockRoomRepository := new(mock_room_repository.MockRoomRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())

	defer func() {
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
	}()

	return mockRoomRepository, mockUoW, &userId, ctx
}

func TestGetRoomQueryHandler(t *testing.T) {
	t.Run("ShouldReturnRoomSuccessWithPlayingSongNotNilAndBoostNil", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setUpMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		userRole := user_role.Host
		now := time.Now().UTC().Truncate(time.Second)
		length := uint16(120)
		roomToBeReturned := &get_room_dao.GetRoomDao{
			Name:                     faker.Name(),
			QrCode:                   []byte(faker.UUIDHyphenated()),
			PlayingSongTitle:         new(faker.Word()),
			PlayingSongAuthor:        new(faker.Name()),
			PlayingSongStartedAtUtc:  new(now),
			PlayingSongLengthSeconds: &length,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           nil,
			BoostCooldownSeconds:     nil,
			SongDaos:                 []get_room_song_dao.GetRoomSongDao{},
		}
		mockRoomRepository.
			On("GetRoomByUserId", ctx, *userId, mock.Anything).
			Return(roomToBeReturned, nil)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)

		expectedQr := base64.RawURLEncoding.EncodeToString(roomToBeReturned.QrCode)
		assert.Equal(t, expectedQr, resp.QrCode)

		// No boost
		assert.Nil(t, resp.BoostData)

		// Playing song mapping
		assert.NotNil(t, resp.PlayingSong)
		assert.Equal(t, *roomToBeReturned.PlayingSongTitle, resp.PlayingSong.Title)
		assert.Equal(t, *roomToBeReturned.PlayingSongAuthor, resp.PlayingSong.Author)
		assert.Equal(t, *roomToBeReturned.PlayingSongStartedAtUtc, resp.PlayingSong.StartedAtUtc)
		assert.Equal(t, *roomToBeReturned.PlayingSongLengthSeconds, resp.PlayingSong.LengthSeconds)

		// No songs
		assert.Equal(t, []get_room_query_response.RoomSongListDto{}, resp.Songs)

		assert.Equal(t, roomToBeReturned.UserRole, resp.UserRole)
	})
	t.Run("ShouldReturnRoomSuccessWithBoostNotNil", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setUpMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		userRole := user_role.Member
		now := time.Now().UTC().Truncate(time.Second)
		boostUsed := now.Add(-5 * time.Minute)
		boostCooldown := uint16(30)
		length := uint16(200)
		song1ID := uuid.New()
		song2ID := uuid.New()
		roomToBeReturned := &get_room_dao.GetRoomDao{
			Name:                     faker.Name(),
			QrCode:                   []byte(faker.UUIDHyphenated()),
			PlayingSongTitle:         new(faker.Word()),
			PlayingSongAuthor:        new(faker.Name()),
			PlayingSongStartedAtUtc:  new(now),
			PlayingSongLengthSeconds: &length,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           &boostUsed,
			BoostCooldownSeconds:     &boostCooldown,
			SongDaos: []get_room_song_dao.GetRoomSongDao{
				{
					Id:            song1ID,
					Title:         faker.Word(),
					Author:        faker.Name(),
					AddedBy:       faker.LastName(),
					State:         enqueued_song_state.Playing.String(),
					Votes:         uint8(3),
					AlbumCoverUrl: faker.URL(),
					VoteStatus:    vote_status.Upvoted.String(),
				},
				{
					Id:            song2ID,
					Title:         faker.Word(),
					Author:        faker.Name(),
					AddedBy:       faker.LastName(),
					State:         enqueued_song_state.Enqueued.String(),
					Votes:         uint8(1),
					AlbumCoverUrl: faker.URL(),
					VoteStatus:    vote_status.NotVoted.String(),
				},
			},
		}
		mockRoomRepository.
			On("GetRoomByUserId", ctx, *userId, mock.Anything).
			Return(roomToBeReturned, nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)

		expectedQr := base64.RawURLEncoding.EncodeToString(roomToBeReturned.QrCode)
		assert.Equal(t, expectedQr, resp.QrCode)

		// Boost mapping
		assert.NotNil(t, resp.BoostData)
		assert.Equal(t, boostUsed, resp.BoostData.BoostUsedAtUtc)
		assert.Equal(t, boostCooldown, resp.BoostData.BoostCooldownSeconds)

		// Playing song mapping
		assert.NotNil(t, resp.PlayingSong)
		assert.Equal(t, *roomToBeReturned.PlayingSongTitle, resp.PlayingSong.Title)
		assert.Equal(t, *roomToBeReturned.PlayingSongAuthor, resp.PlayingSong.Author)
		assert.Equal(t, *roomToBeReturned.PlayingSongStartedAtUtc, resp.PlayingSong.StartedAtUtc)
		assert.Equal(t, *roomToBeReturned.PlayingSongLengthSeconds, resp.PlayingSong.LengthSeconds)

		// Songs mapping
		if assert.Len(t, resp.Songs, 2) {
			s1 := resp.Songs[0]
			assert.Equal(t, song1ID, s1.Id)
			assert.Equal(t, roomToBeReturned.SongDaos[0].Title, s1.Title)
			assert.Equal(t, roomToBeReturned.SongDaos[0].Author, s1.Author)
			assert.Equal(t, roomToBeReturned.SongDaos[0].AddedBy, s1.AddedBy)
			assert.Equal(t, roomToBeReturned.SongDaos[0].Votes, s1.Votes)
			assert.Equal(t, enqueued_song_state.Playing.String(), s1.State)
			assert.Equal(t, vote_status.Upvoted.String(), s1.VoteStatus)

			s2 := resp.Songs[1]
			assert.Equal(t, song2ID, s2.Id)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Title, s2.Title)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Author, s2.Author)
			assert.Equal(t, roomToBeReturned.SongDaos[1].AddedBy, s2.AddedBy)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Votes, s2.Votes)
			assert.Equal(t, enqueued_song_state.Enqueued.String(), s2.State)
			assert.Equal(t, vote_status.NotVoted.String(), s2.VoteStatus)
		}

		assert.Equal(t, roomToBeReturned.UserRole, resp.UserRole)
	})
	t.Run("ShouldReturnRoomSuccess", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setUpMocks(t)

		mockUoW.On("GetQueryer").Return(nil)
		userRole := user_role.Host
		roomToBeReturned := &get_room_dao.GetRoomDao{
			Name:                     "Test Room Name",
			QrCode:                   []byte("Test QrCode"),
			PlayingSongTitle:         nil,
			PlayingSongAuthor:        nil,
			PlayingSongStartedAtUtc:  nil,
			PlayingSongLengthSeconds: nil,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           nil,
			BoostCooldownSeconds:     nil,
			SongDaos:                 []get_room_song_dao.GetRoomSongDao{},
		}
		mockRoomRepository.On("GetRoomByUserId", ctx, *userId, mock.Anything).Return(roomToBeReturned, nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		roomResponse, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		assert.Equal(t, roomToBeReturned.Name, roomResponse.Name)
		assert.Equal(t, "VGVzdCBRckNvZGU", roomResponse.QrCode)
		assert.Nil(t, roomResponse.BoostData)
		assert.Nil(t, roomResponse.PlayingSong)
		assert.Equal(t, roomToBeReturned.UserRole, roomResponse.UserRole)
		assert.Equal(t, []get_room_query_response.RoomSongListDto{}, roomResponse.Songs)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepository, mockUoW, _, _ := setUpMocks(t)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		// Act
		resp, err := handler.Handle(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 0)
		assert.Equal(t, application_helpers.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository, mockUoW, userId, ctx := setUpMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)
		repoErr := errors.New("database error")
		errorCode := "GetRoomQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, *userId, mock.Anything).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_error.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
	})
}
