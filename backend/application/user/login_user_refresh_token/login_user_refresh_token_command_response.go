package login_user_refresh_token

type LoginUserRefreshTokenCommandResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
