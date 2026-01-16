package constants

import "time"

const (
	RoomPlayStateCookieExpirationTime    = time.Minute * 5
	AuthBasePath                         = "/auth"
	LogoutPath                           = "/logout"
	RefreshTokenPath                     = "/refresh-token"
	RoomPlayStateCookieName              = "roomPlay-state"
	RoomPlayAccessTokenCookieName        = "roomplay-session-at"
	RoomPlayRefreshTokenCookieName       = "roomplay-session-rt"
	RoomPlayDeviceIdCookieName           = "roomplay-device-id"
	RoomPlayDeviceTypeCookieName         = "roomplay-device-type"
	RoomPlayDeviceIdCookieExpirationTime = 24 * time.Hour * 365 // a year
	ApiBasePath                          = "/api/v1"
)
