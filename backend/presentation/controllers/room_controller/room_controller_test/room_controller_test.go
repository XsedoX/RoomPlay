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
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/client_message/client_message_envelope"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/event_handlers/song_enqueued_websocket_event"
	"github.com/XsedoX/RoomPlay/infrastructure/websocket/websocket_action"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/authentication_mocks/mock_encrypter"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_music_data_provider_service"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	tests_initializer.InjectEncrypter = func() *mock_encrypter.MockEncrypter {
		mockEncrypter := mock_encrypter.MockEncrypter{}
		mockEncrypter.On(
			"NewRandomString").Return("random_string_for_testingstststst") // NOTE: Must be 32 characters long.
		return &mockEncrypter
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
	require.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[get_room_query_response.GetRoomQueryResponse]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	require.Equal(t, roomName, responseWrapper.Data.Name)
	require.Equal(t, roomToTest.QrCode(), responseWrapper.Data.QrCode)
	require.Equal(t, "host", responseWrapper.Data.UserRole)

	require.NotNil(t, responseWrapper.Data.PlayingSong)
	playingSong := roomToTest.PlayingSong()
	require.Equal(t, playingSong.SongData().Title(), responseWrapper.Data.PlayingSong.Title)
	require.Equal(t, playingSong.SongData().Author(), responseWrapper.Data.PlayingSong.Author)
	require.Equal(t, playingSong.SongData().LengthSeconds(), responseWrapper.Data.PlayingSong.LengthSeconds)

	require.Empty(t, responseWrapper.Data.Songs)
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

	require.Equal(t, http.StatusCreated, w.Code)
	var roomExists bool
	err = txx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	require.NoError(t, err)
	require.Equal(t, true, roomExists)
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

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var roomExists bool
	err = txx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	require.Error(t, sql.ErrNoRows, err)
}

func TestCheckUserRoomMembershipSuccess(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodGet, constants.ApiBasePath+room_controller.RoomBasePath+room_controller.RoomMembershipBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var responseWrapper test_helpers.TestResponseWrapper[bool]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)
	require.True(t, responseWrapper.Data)
}

func TestLeaveRoomSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(
		http.MethodDelete,
		constants.ApiBasePath+room_controller.RoomBasePath,
		nil,
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var isUserInRoom bool
	_ = txx.Get(&isUserInRoom,
		"SELECT EXISTS (SELECT 1 FROM users_room_data WHERE user_id = $1)::text;",
		tests_initializer.InjectedUser.Id())
	require.Equal(t, false, isUserInRoom)
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
	require.Equal(t, true, isUserInRoom)
}

func TestEnqueueSongSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	conn := connectWebSocket(t)

	command := enqueue_song_command.EnqueueSongCommand{
		SongExternalId: seeder.ExternalSongData.Songs[0].VideoId,
	}
	payload, err := json.Marshal(command)
	require.NoError(t, err)
	envelope := client_message_envelope.ClientMessageEnvelope{
		Payload:    payload,
		ActionName: websocket_action.EnqueueSongClientMessageActionName,
	}

	err = conn.WriteJSON(envelope)
	require.NoError(t, err)

	responseBytes := test_helpers.ReadMessageForAction(
		t,
		conn,
		websocket_action.EnqueuedSongPatchActionName,
	)

	var responseSuccess test_helpers.WebSocketPatchTestResponseWrapper
	err = json.Unmarshal(responseBytes, &responseSuccess)
	require.NoError(t, err)
	require.Len(t, responseSuccess.Payload.Data, 1)
	jsonPatch := responseSuccess.Payload.Data[0]
	require.Equal(t, "add", jsonPatch.Op)

	var songData song_enqueued_websocket_event.SongEnqueuedWebsocketEventResponse
	err = json.Unmarshal(jsonPatch.Value, &songData)
	require.NoError(t, err)
	require.Equal(t, seeder.ExternalSongDataById.Title, songData.Title)
	require.Equal(t, seeder.ExternalSongDataById.Author, songData.Author)
	require.Equal(t, seeder.ExternalSongDataById.AlbumCoverUrl, songData.AlbumCoverUrl)
	require.Equal(t, int8(0), songData.Votes)
	require.Equal(t, tests_initializer.InjectedUser.FullName().String(), songData.AddedBy)
	require.Equal(t, websocket_action.EnqueuedSongPatchActionName, responseSuccess.ActionName)

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

func connectWebSocket(t *testing.T) *websocket.Conn {
	t.Helper()
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	req := httptest.NewRequest(http.MethodGet,
		constants.ApiBasePath+
			room_controller.RoomBasePath+
			room_controller.WebSocketTicketPath,
		nil,
	)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	type responseType struct {
		Token string `json:"token"`
	}

	var responseWrapper test_helpers.TestResponseWrapper[responseType]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)
	ticket := responseWrapper.Data
	require.NotEmpty(t, ticket)

	require.NotEmpty(t, ticket.Token)

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	wsUrl := "ws" +
		strings.TrimPrefix(server.URL, "http") +
		constants.ApiBasePath +
		room_controller.RoomBasePath +
		room_controller.WebSocketUpgradePath +
		"?token=" + ticket.Token
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}
