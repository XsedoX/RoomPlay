package room_controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/XsedoX/RoomPlay/application/dtos/music_data_response_dto"
	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/application/room/enqueue_song/enqueue_song_command"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_envelope"
	"github.com/XsedoX/RoomPlay/infrastructure/event_handlers/song_enqueued_websocket_event"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	songToReturn := seeder.ExternalSongDataById
	tests_initializer.InjectMusicDataService = func() *mock_music_data_provider_service.MockMusicDataProviderService {
		mockMusicDataService := mock_music_data_provider_service.MockMusicDataProviderService{}
		mockMusicDataService.On(
			"GetSongById",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&music_data_response_dto.SongDataByIdResponseDto{
			VideoId:       songToReturn.VideoId,
			Title:         songToReturn.Title,
			Author:        songToReturn.Author,
			AlbumCoverUrl: songToReturn.AlbumCoverUrl,
			MusicProvider: songToReturn.MusicProvider,
			Isrc:          songToReturn.Isrc,
			LengthSeconds: songToReturn.LengthSeconds,
		}, nil)
		return &mockMusicDataService
	}
	tests_initializer.InitializeApiServer(m)
}

func TestGetRoomSuccess(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	roomToTest := seeder.SeedData.Rooms[1]
	roomName := roomToTest.Name()

	// Perform Request
	req := httptest.NewRequest(http.MethodGet, constants.ApiBasePath+room_controller.RoomBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[get_room_query_response.GetRoomQueryResponse]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	assert.Equal(t, roomName, responseWrapper.Data.Name)
	assert.Equal(t, roomToTest.QrCode(), responseWrapper.Data.QrCode)
	assert.Equal(t, "host", responseWrapper.Data.UserRole)

	require.NotNil(t, responseWrapper.Data.PlayingSong)
	playingSong := roomToTest.PlayingSong()
	assert.Equal(t, playingSong.SongData().Title(), responseWrapper.Data.PlayingSong.Title)
	assert.Equal(t, playingSong.SongData().Author(), responseWrapper.Data.PlayingSong.Author)
	assert.Equal(t, playingSong.SongData().LengthSeconds(), responseWrapper.Data.PlayingSong.LengthSeconds)

	assert.Empty(t, responseWrapper.Data.Songs)
}

func TestCreateRoomSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	command := create_room_command.CreateRoomCommand{
		RoomName:           "New Test Room",
		RoomPassword:       "password123",
		RepeatRoomPassword: "password123",
	}

	body, err := json.Marshal(command)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+room_controller.RoomBasePath, bytes.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var roomExists bool
	err = txx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	assert.NoError(t, err)
	assert.Equal(t, true, roomExists)
}

func TestCreateRoomValidationFailure(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, false)
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	command := create_room_command.CreateRoomCommand{
		RoomName:           "",
		RoomPassword:       "short",
		RepeatRoomPassword: "mismatch",
	}

	body, err := json.Marshal(command)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		constants.ApiBasePath+room_controller.RoomBasePath,
		bytes.NewReader(body),
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var roomExists bool
	err = txx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	assert.Error(t, sql.ErrNoRows, err)
}

func TestCheckUserRoomMembershipSuccess(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodGet, constants.ApiBasePath+room_controller.RoomBasePath+room_controller.RoomMembershipBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseWrapper test_helpers.TestResponseWrapper[bool]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)
	assert.True(t, responseWrapper.Data)
}

func TestLeaveRoomSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodDelete, constants.ApiBasePath+room_controller.RoomBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var isUserInRoom bool
	_ = txx.Get(&isUserInRoom,
		"SELECT EXISTS (SELECT 1 FROM users_room_data WHERE user_id = $1)::text;",
		tests_initializer.InjectedUser.Id())
	assert.Equal(t, false, isUserInRoom)
}

func TestJoinRoomSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	command := join_room_password_command.JoinRoomPasswordCommand{
		RoomName:     seeder.SeedData.Rooms[0].Name(),
		RoomPassword: string(seeder.SeedData.Rooms[0].Password()),
	}
	body, err := json.Marshal(command)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut,
		constants.ApiBasePath+room_controller.RoomBasePath+room_controller.JoinRoomPasswordPath,
		bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
	var isUserInRoom bool
	_ = txx.Get(&isUserInRoom,
		"SELECT EXISTS (SELECT 1 FROM users_room_data WHERE user_id = $1)::text;",
		tests_initializer.InjectedUser.Id())
	assert.Equal(t, true, isUserInRoom)
}

func TestRoomDataWelcomeMessageSuccess(t *testing.T) {
	_, _ = tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	server := httptest.NewServer(r)
	defer server.Close()

	wsUrl := "ws" +
		strings.TrimPrefix(server.URL, "http") +
		constants.ApiBasePath +
		room_controller.RoomBasePath +
		room_controller.WebSocketUpgradePath
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	require.NoError(t, err)
	defer conn.Close()

	responseBytes := test_helpers.ReadMessageForAction(
		t,
		conn,
		websocket_action.GetRoomDataActionName,
	)

	var responseSuccess test_helpers.WebSocketTestResponseWrapper
	err = json.Unmarshal(responseBytes, &responseSuccess)
	require.Equal(t, websocket_action.GetRoomDataActionName, responseSuccess.ActionName)
	require.NoError(t, err)

	var roomData get_room_query_response.GetRoomQueryResponse
	err = json.Unmarshal(responseSuccess.Data, &roomData)
	require.NoError(t, err)
	require.Equal(t, tests_initializer.InjectedUserRoom.Name(), roomData.Name)
	require.Equal(t, tests_initializer.InjectedUserRoom.QrCode(), roomData.QrCode)
	require.Equal(t, len(tests_initializer.InjectedUserRoom.EnqueuedSongs()), len(roomData.Songs))
	require.Equal(t, *tests_initializer.InjectedUser.Role().String(), roomData.UserRole)
	require.Equal(t, tests_initializer.InjectedUserRoom.PlayingSong().SongData().Title(), roomData.PlayingSong.Title)
}

func TestEnqueueSongSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	server := httptest.NewServer(r)
	defer server.Close()

	wsUrl := "ws" +
		strings.TrimPrefix(server.URL, "http") +
		constants.ApiBasePath +
		room_controller.RoomBasePath +
		room_controller.WebSocketUpgradePath
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	require.NoError(t, err)
	defer conn.Close()

	command := enqueue_song_command.EnqueueSongCommand{
		SongExternalId: seeder.ExternalSongData.Songs[0].VideoId,
	}
	payload, err := json.Marshal(command)
	require.NoError(t, err)
	envelope := client_message_envelope.ClientMessageEnvelope{
		Payload:    payload,
		ActionName: websocket_action.SongEnqueuedClientMessageActionName,
	}

	err = conn.WriteJSON(envelope)
	require.NoError(t, err)

	responseBytes := test_helpers.ReadMessageForAction(
		t,
		conn,
		websocket_action.EnqueuedSongsPatchActionName,
	)

	var responseSuccess test_helpers.WebSocketPatchTestResponseWrapper
	err = json.Unmarshal(responseBytes, &responseSuccess)
	require.NoError(t, err)
	require.Len(t, responseSuccess.Data, 1)
	jsonPatch := responseSuccess.Data[0]
	require.Equal(t, "add", jsonPatch.Op)

	var songData song_enqueued_websocket_event.SongEnqueuedWebsocketEventResponse
	err = json.Unmarshal(jsonPatch.Value, &songData)
	require.NoError(t, err)
	assert.Equal(t, seeder.ExternalSongDataById.Title, songData.Title)
	assert.Equal(t, seeder.ExternalSongDataById.Author, songData.Author)
	assert.Equal(t, seeder.ExternalSongDataById.AlbumCoverUrl, songData.AlbumCoverUrl)
	assert.Equal(t, int8(0), songData.Votes)
	assert.Equal(t, tests_initializer.InjectedUser.FullName().String(), songData.AddedBy)
	assert.Equal(t, websocket_action.EnqueuedSongsPatchActionName, responseSuccess.ActionName)

	var enqueuedSongId uuid.UUID
	_ = txx.Get(&enqueuedSongId,
		`
				SELECT es.id FROM enqueued_songs es 
		   JOIN songs_external_data sed ON es.song_id = sed.song_id 
		   WHERE sed.external_id = $1 AND es.room_id = $2;
		`,
		command.SongExternalId,
		seeder.SeedData.Users[0].RoomId().ToUuid(),
	)
	require.Equal(t, enqueuedSongId.String(), jsonPatch.Path)
}
