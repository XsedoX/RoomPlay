package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"xsedox.com/main/application"
	"xsedox.com/main/application/user"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/entities"
	"xsedox.com/main/presentation/presentationErrors"
)

const (
	roomplayKeyCookieExpirationTime = time.Minute * 5
	googleCallbackEndpoint          = "/api/v1/auth/google/callback"
	stateCookieName                 = "roomPlay-state"
)

//TODO save where the user started logging and return to the same url

type OidcHandler struct {
	loginUserCommandHandler application.ICommandHandler[*user.LoginCommand]
}

func NewOidcHandler(loginUserCommandHandler application.ICommandHandler[*user.LoginCommand]) *OidcHandler {
	return &OidcHandler{
		loginUserCommandHandler: loginUserCommandHandler,
	}
}

func (handler *OidcHandler) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	clearStateCookie(w)
	expiresAt := time.Now().Add(roomplayKeyCookieExpirationTime)
	state := uuid.NewString()
	http.SetCookie(w, &http.Cookie{
		Name:     "deviceFingerprint",
		Value:    state,
		Expires:  expiresAt,
		Path:     googleCallbackEndpoint,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	googleGetUrl, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Invalid redirect URL", http.StatusForbidden)
		return
	}

	parameters := url.Values{}
	parameters.Add("response_type", "code")
	parameters.Add("client_id", config.Config.Authentication.ClientId)
	parameters.Add("scope", config.Config.Scopes)
	parameters.Add("access_type", "offline")
	parameters.Add("redirect_uri", getBackendUrl())
	parameters.Add("state", state)
	googleGetUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, googleGetUrl.String(), http.StatusTemporaryRedirect)
}

func (handler *OidcHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie(stateCookieName)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Invalid state", http.StatusForbidden)
		return
	}
	stateString := state.Value

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Invalid params", http.StatusForbidden)
		return
	}
	if params.Get("state") != stateString {
		presentationErrors.WriteJsonFailure(w, "State mismatch", http.StatusForbidden)
		return
	}
	code := params.Get("code")
	if code == "" {
		presentationErrors.WriteJsonFailure(w, "Invalid code", http.StatusForbidden)
		return
	}
	tokenURL := "https://oauth2.googleapis.com/token"

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", config.Config.Authentication.ClientId)
	form.Add("client_secret", config.Config.Authentication.ClientSecret)
	form.Add("redirect_uri", getBackendUrl())

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Failed to create token request", http.StatusForbidden)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Failed to call token endpoint", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		presentationErrors.WriteJsonFailure(w, "Token endpoint returned non-200", http.StatusForbidden)
		return
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
		IdToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		presentationErrors.WriteJsonFailure(w, "Failed to decode token response", http.StatusForbidden)
		return
	}
	type googleApiClaims struct {
		jwt.RegisteredClaims
		GivenName  string `json:"given_name" validate:"required"`
		FamilyName string `json:"family_name" validate:"required"`
	}
	token, err := jwt.ParseWithClaims(tokenResp.IdToken, &googleApiClaims{}, nil)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Invalid Auth", http.StatusUnauthorized)
		return
	}
	claims, ok := token.Claims.(*googleApiClaims)
	if !ok || !token.Valid {
		presentationErrors.WriteJsonFailure(w, "Invalid Auth", http.StatusUnauthorized)
		return
	}
	fingerprintHeader := r.Header.Get("X-DeviceFingerprint")
	if fingerprintHeader == "" {
		presentationErrors.WriteJsonFailure(w, "Missing fingerprint", http.StatusBadRequest)
		return
	}
	loginUserCommand := user.LoginCommand{
		Name: claims.GivenName,
		Device: user.DeviceDto{
			Fingerprint: fingerprintHeader,
			DeviceType:  entities.COMPUTER,
		},
		ExternalId: claims.Subject,
		Surname:    claims.FamilyName,
	}

	err = handler.loginUserCommandHandler.Handle(r.Context(), &loginUserCommand)
	if err != nil {
		presentationErrors.WriteJsonFailure(w, "Failed to login", http.StatusForbidden)
		return
	}

	http.Redirect(w, r, config.Config.Authentication.ClientOrigin+"/signin-oidc", http.StatusPermanentRedirect)
}

func clearStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   stateCookieName,
		MaxAge: -1,
	})
}

func getBackendUrl() string {
	return config.Config.Server.Host + ":" + config.Config.Server.Port +
		googleCallbackEndpoint
}
