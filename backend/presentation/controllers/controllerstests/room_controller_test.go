package controllerstests

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/room/create_room"
	"github.com/XsedoX/RoomPlay/application/room/get_room"
	"github.com/XsedoX/RoomPlay/presentation/controllers"
	"github.com/XsedoX/RoomPlay/presentation/helpers"
	"github.com/XsedoX/RoomPlay/test_helpers"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoomSuccess(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()
	roomName := integration_tests.SeedData.Rooms[1].Name()

	// Perform Request
	req := httptest.NewRequest(http.MethodGet, helpers.ApiBasePath+controllers.RoomBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[get_room.GetRoomQueryResponse]
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
	testServer := integration_tests.TestServer
	r := testServer.Router()

	command := create_room.CreateRoomCommand{
		RoomName:           "New Test Room",
		RoomPassword:       "password123",
		RepeatRoomPassword: "password123",
	}

	body, err := json.Marshal(command)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, helpers.ApiBasePath+controllers.RoomBasePath, bytes.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var roomExists bool
	dbx := integration_tests.PgContainer.DB
	err = dbx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	assert.NoError(t, err)
	assert.Equal(t, true, roomExists)
}

func TestCreateRoomValidationFailure(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()

	command := create_room.CreateRoomCommand{
		RoomName:           "",
		RoomPassword:       "short",
		RepeatRoomPassword: "mismatch",
	}

	body, err := json.Marshal(command)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, helpers.ApiBasePath+controllers.RoomBasePath, bytes.NewReader(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var roomExists bool
	dbx := integration_tests.PgContainer.DB
	err = dbx.Get(&roomExists, "SELECT EXISTS (SELECT * FROM rooms WHERE name = $1)::text;", command.RoomName)
	assert.Error(t, sql.ErrNoRows, err)
}

func TestCheckUserRoomMembershipSuccess(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodGet, helpers.ApiBasePath+controllers.RoomBasePath+controllers.RoomMembershipBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[bool]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	assert.True(t, responseWrapper.Data)
}

func TestLeaveRoomSuccess(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodDelete, helpers.ApiBasePath+controllers.RoomBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	dbx := integration_tests.PgContainer.DB
	var isUserInRoom bool
	_ = dbx.Get(&isUserInRoom,
		"SELECT EXISTS (SELECT 1 FROM users_room_data WHERE user_id = $1)::text;",
		integration_tests.InjectedUser.Id())
	assert.Equal(t, false, isUserInRoom)
}
