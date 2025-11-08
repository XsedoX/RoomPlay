package helpers

import "time"

const (
	RoomplayStateCookieExpirationTime    = time.Minute * 5
	RoomplayStateCookieName              = "roomPlay-state"
	RoomplayAccessTokenCookieName        = "roomplay-session-at"
	RoomplayRefreshTokenCookieName       = "roomplay-session-rt"
	RoomPlayDeviceIdCookieName           = "roomplay-device-id"
	RoomPlayDeviceTypeCookieName         = "roomplay-device-type"
	RoomPlayDeviceIdCookieExpirationTime = 24 * time.Hour * 365 // a year
)
