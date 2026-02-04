package room_controller_test

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/room/create_room/create_room_command"
	"github.com/XsedoX/RoomPlay/application/room/get_room/get_room_query_response"
	"github.com/XsedoX/RoomPlay/application/room/join_room_password/join_room_password_command"
	"github.com/XsedoX/RoomPlay/presentation/controllers/room_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeApiServer(m)
}

func TestGetRoomSuccess(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	roomName := seeder.SeedData.Rooms[1].Name()

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
	assert.Equal(t, base64.RawURLEncoding.EncodeToString([]byte("qrCode2")), responseWrapper.Data.QrCode)
	assert.Equal(t, "host", responseWrapper.Data.UserRole)

	require.NotNil(t, responseWrapper.Data.PlayingSong)
	assert.Equal(t, "title2", responseWrapper.Data.PlayingSong.Title)
	assert.Equal(t, "author2", responseWrapper.Data.PlayingSong.Author)
	assert.Equal(t, uint16(349), responseWrapper.Data.PlayingSong.LengthSeconds)

	assert.NotEmpty(t, responseWrapper.Data.Songs)
	assert.Len(t, responseWrapper.Data.Songs, 5)
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

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+room_controller.RoomBasePath, bytes.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
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
		RoomPassword: seeder.SeedData.Rooms[0].Password(),
	}
	body, err := json.Marshal(command)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut,
		constants.ApiBasePath+room_controller.RoomBasePath+room_controller.JoinRoomPasswordPath,
		bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	result := w.Result()
	bodyBytes, _ := io.ReadAll(result.Body)
	print(string(bodyBytes))
	require.Equal(t, http.StatusNoContent, w.Code)
	var isUserInRoom bool
	_ = txx.Get(&isUserInRoom,
		"SELECT EXISTS (SELECT 1 FROM users_room_data WHERE user_id = $1)::text;",
		tests_initializer.InjectedUser.Id())
	assert.Equal(t, true, isUserInRoom)
}
