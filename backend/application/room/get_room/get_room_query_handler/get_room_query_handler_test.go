package get_room_query_handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_helpers"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/daos/get_room_song_dao"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/enqueued_song_state"
	"github.com/XsedoX/RoomPlay/domain/room/enqueued_song/vote_status"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_role"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_room_repository"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/persistance_mocks/mock_unit_of_work"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setUpMocks(t *testing.T) (*mock_room_repository.MockRoomRepository,
	*mock_unit_of_work.MockUnitOfWork,
	*mock_encrypter.MockEncrypter,
	*user_id.UserId,
	context.Context,
) {
	mockRoomRepository := new(mock_room_repository.MockRoomRepository)
	mockUoW := new(mock_unit_of_work.MockUnitOfWork)
	userId, ctx := test_helpers.AddUserIdToContext(context.Background())
	mockEncrypter := new(mock_encrypter.MockEncrypter)

	defer func() {
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockEncrypter.AssertExpectations(t)
	}()

	return mockRoomRepository, mockUoW, mockEncrypter, &userId, ctx
}

func TestGetRoomQueryHandler(t *testing.T) {
	t.Run("ShouldReturnRoomSuccessWithPlayingSongNotNilAndBoostNil", func(t *testing.T) {
		mockRoomRepository, mockUoW, mockEncrypter, userId, ctx := setUpMocks(t)

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

		mockEncrypter.On("Decrypt", roomToBeReturned.QrCode).Return(string(roomToBeReturned.QrCode), nil)
		mockRoomRepository.
			On("GetRoomByUserId", ctx, *userId, mock.Anything).
			Return(roomToBeReturned, nil)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository, mockEncrypter)

		resp, err := handler.Handle(ctx)

		require.NoError(t, err)
		require.NotNil(t, resp)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)

		require.Equal(t, string(roomToBeReturned.QrCode), resp.QrCode)

		// No boost
		require.Nil(t, resp.BoostData)

		// Playing song mapping
		require.NotNil(t, resp.PlayingSong)
		require.Equal(t, *roomToBeReturned.PlayingSongTitle, resp.PlayingSong.Title)
		require.Equal(t, *roomToBeReturned.PlayingSongAuthor, resp.PlayingSong.Author)
		require.Equal(t, *roomToBeReturned.PlayingSongStartedAtUtc, resp.PlayingSong.StartedAtUtc)
		require.Equal(t, *roomToBeReturned.PlayingSongLengthSeconds, resp.PlayingSong.LengthSeconds)

		// No songs
		require.Equal(t, []get_room_query_response.RoomSongListDto{}, resp.Songs)

		require.Equal(t, roomToBeReturned.UserRole, resp.UserRole)
	})
	t.Run("ShouldReturnRoomSuccessWithBoostNotNil", func(t *testing.T) {
		mockRoomRepository, mockUoW, encrypter, userId, ctx := setUpMocks(t)
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
		encrypter.On("Decrypt", roomToBeReturned.QrCode).Return(string(roomToBeReturned.QrCode), nil)
		mockRoomRepository.
			On("GetRoomByUserId", ctx, *userId, mock.Anything).
			Return(roomToBeReturned, nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository, encrypter)

		resp, err := handler.Handle(ctx)

		require.NoError(t, err)
		require.NotNil(t, resp)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)

		require.Equal(t, string(roomToBeReturned.QrCode), resp.QrCode)

		// Boost mapping
		require.NotNil(t, resp.BoostData)
		require.Equal(t, boostUsed, resp.BoostData.BoostUsedAtUtc)
		require.Equal(t, boostCooldown, resp.BoostData.BoostCooldownSeconds)

		// Playing song mapping
		require.NotNil(t, resp.PlayingSong)
		require.Equal(t, *roomToBeReturned.PlayingSongTitle, resp.PlayingSong.Title)
		require.Equal(t, *roomToBeReturned.PlayingSongAuthor, resp.PlayingSong.Author)
		require.Equal(t, *roomToBeReturned.PlayingSongStartedAtUtc, resp.PlayingSong.StartedAtUtc)
		require.Equal(t, *roomToBeReturned.PlayingSongLengthSeconds, resp.PlayingSong.LengthSeconds)

		// Songs mapping
		require.Len(t, resp.Songs, 2)
		s1 := resp.Songs[0]
		require.Equal(t, song1ID, s1.Id)
		require.Equal(t, roomToBeReturned.SongDaos[0].Title, s1.Title)
		require.Equal(t, roomToBeReturned.SongDaos[0].Author, s1.Author)
		require.Equal(t, roomToBeReturned.SongDaos[0].AddedBy, s1.AddedBy)
		require.Equal(t, roomToBeReturned.SongDaos[0].Votes, s1.Votes)
		require.Equal(t, enqueued_song_state.Playing.String(), s1.State)
		require.Equal(t, vote_status.Upvoted.String(), s1.VoteStatus)

		s2 := resp.Songs[1]
		require.Equal(t, song2ID, s2.Id)
		require.Equal(t, roomToBeReturned.SongDaos[1].Title, s2.Title)
		require.Equal(t, roomToBeReturned.SongDaos[1].Author, s2.Author)
		require.Equal(t, roomToBeReturned.SongDaos[1].AddedBy, s2.AddedBy)
		require.Equal(t, roomToBeReturned.SongDaos[1].Votes, s2.Votes)
		require.Equal(t, enqueued_song_state.Enqueued.String(), s2.State)
		require.Equal(t, vote_status.NotVoted.String(), s2.VoteStatus)

		require.Equal(t, roomToBeReturned.UserRole, resp.UserRole)
	})
	t.Run("ShouldReturnRoomSuccess", func(t *testing.T) {
		mockRoomRepository, mockUoW, encrypter, userId, ctx := setUpMocks(t)

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
		encrypter.On("Decrypt", roomToBeReturned.QrCode).Return(string(roomToBeReturned.QrCode), nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository, encrypter)

		roomResponse, err := handler.Handle(ctx)

		require.NoError(t, err)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
		require.Equal(t, roomToBeReturned.Name, roomResponse.Name)
		require.Equal(t, string(roomToBeReturned.QrCode), roomResponse.QrCode)
		require.Nil(t, roomResponse.BoostData)
		require.Nil(t, roomResponse.PlayingSong)
		require.Equal(t, roomToBeReturned.UserRole, roomResponse.UserRole)
		require.Equal(t, []get_room_query_response.RoomSongListDto{}, roomResponse.Songs)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepository, mockUoW, encrypter, _, _ := setUpMocks(t)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository, encrypter)

		// Act
		resp, err := handler.Handle(context.Background())

		// Assert
		require.Error(t, err)
		require.Nil(t, resp)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 0)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 0)
		require.Equal(t, application_helpers.NewMissingUserIdInContextError("GetRoomQueryHandler.Handle"), err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository, mockUoW, encrypter, userId, ctx := setUpMocks(t)
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository, encrypter)
		repoErr := errors.New("database error")
		errorCode := "GetRoomQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, *userId, mock.Anything).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		require.Error(t, err)
		require.Nil(t, resp)
		var customErr *application_error.ApplicationError
		require.True(t, errors.As(err, &customErr))
		require.Equal(t, errorCode, customErr.Code)
		require.ErrorIs(t, customErr.Err, repoErr)
		mockUoW.AssertNumberOfCalls(t, "GetQueryer", 1)
		mockRoomRepository.AssertNumberOfCalls(t, "GetRoomByUserId", 1)
	})
}
