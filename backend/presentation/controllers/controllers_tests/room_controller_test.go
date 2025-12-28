package controllers_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xsedox.com/main/application/room/get_room"
	"xsedox.com/main/test_helpers/integration_tests"
)

func TestGetRoom(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()
	// Middleware to inject user ID
	roomName := integration_tests.SeedData.Rooms[1].Name()

	// Perform Request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/room", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper struct {
		Data get_room.GetRoomQueryResponse `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	assert.Equal(t, roomName, responseWrapper.Data.Name)
}
