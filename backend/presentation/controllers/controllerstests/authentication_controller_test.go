package controllerstests

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/credentials"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/daos"
	"github.com/XsedoX/RoomPlay/presentation/helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests"
	othermocks "github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogoutSuccess(t *testing.T) {
	txx, _ := integration_tests.GetTxxAndCtx(t)
	configuration := othermocks.MockConfiguration{}
	testServer := integration_tests.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+constants.AuthBasePath+constants.LogoutPath, nil)
	w := httptest.NewRecorder()
	deviceId := integration_tests.InjectedUser.Devices()[0].Id()
	expiresAt := time.Now().UTC().Add(constants.RoomPlayDeviceIdCookieExpirationTime)
	req.AddCookie(&http.Cookie{
		Name:     constants.RoomPlayDeviceIdCookieName,
		Value:    *deviceId.String(),
		Expires:  expiresAt,
		Path:     configuration.Server().BasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	var tokenFromDb daos.RefreshTokenDao
	err := txx.Get(&tokenFromDb,
		"SELECT * FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2;",
		integration_tests.InjectedUser.Id(), deviceId)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestRefreshTokenSuccess(t *testing.T) {
	txx, _ := integration_tests.GetTxxAndCtx(t)
	configuration := othermocks.MockConfiguration{}
	testServer := integration_tests.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+constants.AuthBasePath+constants.RefreshTokenPath, nil)
	w := httptest.NewRecorder()
	expiresAt := time.Now().UTC().Add(credentials.RefreshTokenExpirationTime)
	deviceId := integration_tests.InjectedUser.Devices()[0].Id()
	encodedRefreshToken := base64.RawURLEncoding.EncodeToString([]byte(integration_tests.SeedData.LoggedInUserRefreshToken.RefreshToken()))
	req.AddCookie(&http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    encodedRefreshToken,
		Expires:  expiresAt,
		Path:     configuration.Server().BasePath + constants.RefreshTokenPath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	var tokenFromDb daos.RefreshTokenDao
	err := txx.Get(&tokenFromDb,
		"SELECT * FROM users_refresh_tokens WHERE user_id = $1 AND device_id = $2;",
		integration_tests.InjectedUser.Id(), deviceId)
	assert.NoError(t, err)
	assert.Greater(t, tokenFromDb.ExpiresAtUtc, time.Now().UTC())

	refreshToken := getCookieByName(w.Result(), constants.RoomPlayRefreshTokenCookieName)
	decodedToken, err := base64.RawURLEncoding.DecodeString(refreshToken.Value)
	assert.NoError(t, err)
	decodedTokenString := string(decodedToken)
	encrypter := authentication.NewEncrypter(&configuration)
	hashedDecodedToken := encrypter.Hash(decodedTokenString)
	assert.Equal(t, tokenFromDb.RefreshToken, hashedDecodedToken)
}

func getCookieByName(resp *http.Response, name string) *http.Cookie {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
