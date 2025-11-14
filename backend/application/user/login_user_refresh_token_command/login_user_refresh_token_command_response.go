package login_user_refresh_token_command

type LoginUserRefreshTokenCommandResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
