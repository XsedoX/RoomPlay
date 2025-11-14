package dtos

type GoogleTokenResponseDto struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
}
