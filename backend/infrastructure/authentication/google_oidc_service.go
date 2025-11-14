package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"xsedox.com/main/application/dtos"
	"xsedox.com/main/config"
)

const (
	googleCallbackEndpoint = "/api/v1/auth/google/callback"
)

type GoogleOidcService struct {
	configuration config.IConfiguration
}

func (g GoogleOidcService) ParseIdToken(idToken string) (*dtos.GoogleIdTokenClaimsDto, error) {
	type googleApiClaims struct {
		jwt.RegisteredClaims
		GivenName  string `json:"given_name" validate:"required"`
		FamilyName string `json:"family_name" validate:"required"`
	}
	googleClaims := googleApiClaims{}
	token, _ := jwt.ParseWithClaims(idToken, &googleClaims, nil)

	claims, ok := token.Claims.(*googleApiClaims)
	if !ok {
		return nil, errors.New("couldn't parse id token")
	}
	return &dtos.GoogleIdTokenClaimsDto{
		GivenName:  claims.GivenName,
		FamilyName: claims.FamilyName,
		Subject:    claims.Subject,
	}, nil
}

func (g GoogleOidcService) GetAccessToken(ctx context.Context, code string) (*dtos.GoogleTokenResponseDto, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", g.configuration.Authentication().ClientId)
	form.Add("client_secret", g.configuration.Authentication().ClientSecret)
	form.Add("redirect_uri", getGoogleCallbackBackendUrl(g.configuration.Server().Host, g.configuration.Server().Port))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var response dtos.GoogleTokenResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (g GoogleOidcService) GenerateOidcUrl(state string) (string, error) {
	googleGetUrl, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		return "", err
	}
	parameters := url.Values{}
	parameters.Add("response_type", "code")
	parameters.Add("client_id", g.configuration.Authentication().ClientId)
	parameters.Add("scope", g.configuration.Authentication().ScopesField)
	parameters.Add("access_type", "offline")
	parameters.Add("redirect_uri", getGoogleCallbackBackendUrl(g.configuration.Server().Host, g.configuration.Server().Port))
	parameters.Add("state", state)
	googleGetUrl.RawQuery = parameters.Encode()
	return googleGetUrl.String(), nil
}

func NewGoogleOidcService(configuration config.IConfiguration) *GoogleOidcService {
	return &GoogleOidcService{configuration: configuration}
}
func getGoogleCallbackBackendUrl(host, port string) string {
	return fmt.Sprintf("%s:%s%s", host, port, googleCallbackEndpoint)
}
