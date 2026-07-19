package authentication_controller_test

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/XsedoX/RoomPlay/domain/internal_credentials"
	"github.com/XsedoX/RoomPlay/infrastructure/authentication/encryper"
	"github.com/XsedoX/RoomPlay/infrastructure/persistance/internal_credentials/internal_credentials_dao"
	"github.com/XsedoX/RoomPlay/presentation/presentation_helpers/constants"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/other_mocks/mock_configuration"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/seeder"
	"github.com/XsedoX/RoomPlay/test_helpers/integration_tests/tests_initializer"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests_initializer.InitializeApiServer(m)
}

func TestLogoutSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+constants.LogoutPath, nil)
	w := httptest.NewRecorder()
	deviceId := tests_initializer.InjectedUser.Devices()[0].Id()
	expiresAt := time.Now().UTC().Add(constants.RoomPlayDeviceIdCookieExpirationTime)
	req.AddCookie(&http.Cookie{
		Name:     constants.RoomPlayDeviceIdCookieName,
		Value:    *deviceId.String(),
		Expires:  expiresAt,
		Path:     constants.ApiBasePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	var tokenFromDb internal_credentials_dao.InternalCredentialsDao
	err := txx.Get(&tokenFromDb,
		"SELECT * FROM users_internal_credentials WHERE user_id = $1 AND device_id = $2;",
		tests_initializer.InjectedUser.Id(), deviceId)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestRefreshTokenSuccess(t *testing.T) {
	txx, _ := tests_initializer.GetTxxAndCtx(t, true)
	configuration := mock_configuration.MockConfiguration{}
	testServer := tests_initializer.TestServer
	r := testServer.Router()

	req := httptest.NewRequest(http.MethodPost, constants.ApiBasePath+constants.RefreshTokenPath, nil)
	w := httptest.NewRecorder()
	expiresAt := time.Now().UTC().Add(internal_credentials.RefreshTokenExpirationTime)
	deviceId := tests_initializer.InjectedUser.Devices()[0].Id()
	encodedRefreshToken := base64.RawURLEncoding.EncodeToString([]byte(seeder.SeedData.InternalCredentials[0].RefreshToken()))
	req.AddCookie(&http.Cookie{
		Name:     constants.RoomPlayRefreshTokenCookieName,
		Value:    encodedRefreshToken,
		Expires:  expiresAt,
		Path:     constants.RefreshTokenCookiePath,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	var tokenFromDb internal_credentials_dao.InternalCredentialsDao
	err := txx.Get(&tokenFromDb,
		"SELECT * FROM users_internal_credentials WHERE user_id = $1 AND device_id = $2;",
		tests_initializer.InjectedUser.Id(), deviceId)
	assert.NoError(t, err)
	assert.Greater(t, tokenFromDb.ExpiresAtUtc, time.Now().UTC())

	refreshToken := getCookieByName(w.Result(), constants.RoomPlayRefreshTokenCookieName)
	decodedToken, err := base64.RawURLEncoding.DecodeString(refreshToken.Value)
	assert.NoError(t, err)
	decodedTokenString := string(decodedToken)
	encrypter := encryper.NewEncrypter(configuration.Authentication().EncryptionKey)
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
