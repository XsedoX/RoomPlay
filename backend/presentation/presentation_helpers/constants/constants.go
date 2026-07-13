package constants

import "time"

const (
	authBasePath   = "/auth"
	logoutBasePath = "/logout"

	RoomPlayStateCookieExpirationTime    = time.Minute * 5
	RefreshTokenBasePath                 = "/refresh-token"
	RoomPlayStateCookieName              = "roomPlay-state"
	RoomPlayAccessTokenCookieName        = "roomplay-session-at"
	RoomPlayRefreshTokenCookieName       = "roomplay-session-rt"
	RoomPlayDeviceIdCookieName           = "roomplay-device-id"
	RoomPlayDeviceTypeCookieName         = "roomplay-device-type"
	RoomPlayDeviceIdCookieExpirationTime = 24 * time.Hour * 365 // a year
	ApiBasePath                          = "/api/v1"
	RefreshTokenPath                     = authBasePath + RefreshTokenBasePath
	LogoutPath                           = authBasePath + logoutBasePath
	RefreshTokenCookiePath               = ApiBasePath + RefreshTokenPath
)
