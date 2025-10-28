package controllers

import (
	"net/http"

	"xsedox.com/main/application/contracts"
)

type AuthenticationController struct {
	refreshToken contracts.IRefreshTokenRepository
}

func NewAuthenticationController(refreshToken contracts.IRefreshTokenRepository) *AuthenticationController {
	return &AuthenticationController{
		refreshToken: refreshToken,
	}
}

func (handler *AuthenticationController) RefreshToken(w http.ResponseWriter, req *http.Request) {

}
