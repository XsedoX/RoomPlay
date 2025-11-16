package get_room_query

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"xsedox.com/main/application"
	"xsedox.com/main/application/custom_errors"
	"xsedox.com/main/application/room/get_room_query/daos"
	"xsedox.com/main/domain/room"
	"xsedox.com/main/domain/user"
	"xsedox.com/main/tests"
	persistance2 "xsedox.com/main/tests/infrustructure/persistance"
)

func TestGetRoomQueryHandler(t *testing.T) {
	t.Run("ShouldReturnRoomSuccessWithPlayingSongNotNilAndBoostNil", func(t *testing.T) {
		mockRoomRepository := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		var f tests.FakeValueProviders
		_ = faker.FakeData(&f)
		userRole := user.Host
		now := time.Now().UTC().Truncate(time.Second)
		length := uint8(120)
		roomToBeReturned := &daos.GetRoomDao{
			Name:                     f.Sentence,
			QrCode:                   []byte(f.UUID),
			PlayingSongTitle:         tests.PtrString(f.Sentence),
			PlayingSongAuthor:        tests.PtrString(f.Name),
			PlayingSongStartedAtUtc:  tests.PtrTime(now),
			PlayingSongLengthSeconds: &length,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           nil,
			BoostCooldownSeconds:     nil,
			SongDaos:                 []daos.GetRoomSongDao{},
		}
		mockRoomRepository.
			On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).
			Return(roomToBeReturned, nil)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)

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
		assert.Equal(t, []RoomSongListDto{}, resp.Songs)

		assert.Equal(t, roomToBeReturned.UserRole, resp.UserRole)
	})
	t.Run("ShouldReturnRoomSuccessWithBoostNotNil", func(t *testing.T) {
		mockRoomRepository := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		userRole := user.Member
		var f1 tests.FakeValueProviders
		err := faker.FakeData(&f1)
		if err != nil {
			fmt.Println(err)
		}
		var f2 tests.FakeValueProviders
		err = faker.FakeData(&f2)
		if err != nil {
			fmt.Println(err)
		}
		now := time.Now().UTC().Truncate(time.Second)
		boostUsed := now.Add(-5 * time.Minute)
		boostCooldown := uint8(30)
		length := uint8(200)
		song1ID := uuid.New()
		song2ID := uuid.New()
		roomToBeReturned := &daos.GetRoomDao{
			Name:                     f1.Sentence,
			QrCode:                   []byte(f1.UUID),
			PlayingSongTitle:         tests.PtrString(f2.Sentence),
			PlayingSongAuthor:        tests.PtrString(f2.Name),
			PlayingSongStartedAtUtc:  tests.PtrTime(now),
			PlayingSongLengthSeconds: &length,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           &boostUsed,
			BoostCooldownSeconds:     &boostCooldown,
			SongDaos: []daos.GetRoomSongDao{
				{
					Id:            song1ID,
					Title:         f1.Word,
					Author:        f1.Name,
					AddedBy:       f1.Word,
					State:         room.Playing,
					Votes:         uint8(3),
					AlbumCoverUrl: f1.Url,
					VoteStatus:    room.Upvoted,
				},
				{
					Id:            song2ID,
					Title:         f2.Word,
					Author:        f2.Name,
					AddedBy:       f2.Word,
					State:         room.Enqueued,
					Votes:         uint8(1),
					AlbumCoverUrl: f2.Url,
					VoteStatus:    room.NotVoted,
				},
			},
		}
		mockRoomRepository.
			On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).
			Return(roomToBeReturned, nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		resp, err := handler.Handle(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)

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
			assert.Equal(t, room.Playing.String(), s1.State)
			assert.Equal(t, room.Upvoted.String(), s1.VoteStatus)

			s2 := resp.Songs[1]
			assert.Equal(t, song2ID, s2.Id)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Title, s2.Title)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Author, s2.Author)
			assert.Equal(t, roomToBeReturned.SongDaos[1].AddedBy, s2.AddedBy)
			assert.Equal(t, roomToBeReturned.SongDaos[1].Votes, s2.Votes)
			assert.Equal(t, room.Enqueued.String(), s2.State)
			assert.Equal(t, room.NotVoted.String(), s2.VoteStatus)
		}

		assert.Equal(t, roomToBeReturned.UserRole, resp.UserRole)

	})
	t.Run("ShouldReturnRoomSuccess", func(t *testing.T) {
		mockRoomRepository := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		userRole := user.Host
		roomToBeReturned := &daos.GetRoomDao{
			Name:                     "Test Room Name",
			QrCode:                   []byte("Test QrCode"),
			PlayingSongTitle:         nil,
			PlayingSongAuthor:        nil,
			PlayingSongStartedAtUtc:  nil,
			PlayingSongLengthSeconds: nil,
			UserRole:                 *userRole.String(),
			BoostUsedAtUtc:           nil,
			BoostCooldownSeconds:     nil,
			SongDaos:                 []daos.GetRoomSongDao{},
		}
		mockRoomRepository.On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).Return(roomToBeReturned, nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)

		roomResponse, err := handler.Handle(ctx)

		assert.NoError(t, err)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		assert.Equal(t, roomToBeReturned.Name, roomResponse.Name)
		assert.Equal(t, "VGVzdCBRckNvZGU", roomResponse.QrCode)
		assert.Nil(t, roomResponse.BoostData)
		assert.Nil(t, roomResponse.PlayingSong)
		assert.Equal(t, roomToBeReturned.UserRole, roomResponse.UserRole)
		assert.Equal(t, []RoomSongListDto{}, roomResponse.Songs)
	})
	t.Run("ShouldReturnErrorWhenUserIdIsMissingFromContext", func(t *testing.T) {
		// Arrange
		mockRoomRepo := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)

		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepo)

		// Act
		resp, err := handler.Handle(context.Background())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRoomRepo.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
		assert.Equal(t, application.NewMissingUserIdInContextError, err)
	})
	t.Run("ShouldReturnErrorWhenRoomRepositoryFails", func(t *testing.T) {
		mockRoomRepository := new(persistance2.MockRoomRepository)
		mockUoW := new(persistance2.MockUnitOfWork)
		userId, ctx := tests.AddUserIdToContext(context.Background())
		mockUoW.On("GetQueryer").Return(nil)
		handler := NewGetRoomQueryHandler(mockUoW, mockRoomRepository)
		repoErr := errors.New("database error")
		errorCode := "GetRoomQueryHandler.GetRoomByUserId"
		mockRoomRepository.On("GetRoomByUserId", ctx, userId, mockUoW.GetQueryer()).Return(nil, repoErr)

		resp, err := handler.Handle(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
		var customErr *custom_errors.CustomError
		assert.True(t, errors.As(err, &customErr))
		assert.Equal(t, errorCode, customErr.Code)
		assert.ErrorIs(t, customErr.Err, repoErr)
		mockRoomRepository.AssertExpectations(t)
		mockUoW.AssertExpectations(t)
	})
}
