package controllerstests

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos"
	"github.com/XsedoX/RoomPlay/presentation/controllers"
	"github.com/XsedoX/RoomPlay/presentation/helpers"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
	othermocks "github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogoutSuccess(t *testing.T) {
	configuration := othermocks.MockConfiguration{}
	testServer := integration_tests.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodPost, helpers.ApiBasePath+controllers.AuthBasePath+controllers.LogoutBasePath, nil)
	w := httptest.NewRecorder()
	deviceId := integration_tests.InjectedUser.Devices()[0].Id()
	expiresAt := time.Now().UTC().Add(helpers.RoomPlayDeviceIdCookieExpirationTime)
	req.AddCookie(&http.Cookie{
		Name:     helpers.RoomPlayDeviceIdCookieName,
		Value:    *deviceId.String(),
		Expires:  expiresAt,
		Path:     configuration.Server().BasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	dbx := integration_tests.PgContainer.DB

	var tokenFromDb daos.RefreshTokenDao
	err := dbx.Get(&tokenFromDb,
		"SELECT * FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2;",
		integration_tests.InjectedUser.Id(), deviceId)
	assert.Equal(t, sql.ErrNoRows, err)
}
