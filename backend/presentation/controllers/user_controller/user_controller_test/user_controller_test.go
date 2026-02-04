package user_controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/user/get_user/get_user_query_response"
	"github.com/XsedoX/RoomPlay/presentation/controllers/user_controller"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/XsedoX/RoomPlay/test_helpers/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeApiServer(m)
}

func TestGetUserData(t *testing.T) {
	testServer := tests_initializer.TestServer
	r := testServer.Router()
	userName := tests_initializer.InjectedUser.FullName().Name()
	userSurname := tests_initializer.InjectedUser.FullName().Surname()

	req := httptest.NewRequest(http.MethodGet, constants.ApiBasePath+user_controller.UserBasePath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseWrapper test_helpers.TestResponseWrapper[get_user_query_response.GetUserQueryResponse]
	err := json.NewDecoder(w.Body).Decode(&responseWrapper)
	require.NoError(t, err)

	assert.Equal(t, userName, responseWrapper.Data.Name)
	assert.Equal(t, userSurname, responseWrapper.Data.Surname)
}
