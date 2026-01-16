package controllerstests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/user/get_user"
	"github.com/XsedoX/RoomPlay/presentation/controllers"
	"github.com/XsedoX/RoomPlay/presentation/helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserData(t *testing.T) {
	testServer := integration_tests.TestServer
	r := testServer.Router()
	userName := integration_tests.InjectedUser.FullName().Name()
	userSurname := integration_tests.InjectedUser.FullName().Surname()

	req := httptest.NewRequest(http.MethodGet, constants.ApiBasePath+controllers.UserBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[get_user.GetUserQueryResponse]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	assert.Equal(t, userName, responseWrapper.Data.Name)
	assert.Equal(t, userSurname, responseWrapper.Data.Surname)
}
