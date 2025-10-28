package authentication

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/shared"
)

const AccessTokenExpirationTime = time.Minute * 5

type JwtProvider struct {
	configuration config.IConfiguration
}

func NewJwtProvider(configuration config.IConfiguration) *JwtProvider {
	return &JwtProvider{
		configuration: configuration,
	}
}

func (jwtProvider *JwtProvider) GenerateToken(userId shared.UserId) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid.UUID(userId).String(),
		"aud": jwtProvider.configuration.Authentication().AudienceField,
		"iss": jwtProvider.configuration.Authentication().Issuer,
		"exp": jwt.NewNumericDate(time.Now().Add(AccessTokenExpirationTime).UTC()),
		"iat": jwt.NewNumericDate(time.Now().UTC()),
	})
	tokenString, err := token.SignedString([]byte(jwtProvider.configuration.Authentication().JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (jwtProvider *JwtProvider) ValidateTokenAndGetUserId(tokenString string) (*shared.UserId, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtProvider.configuration.Authentication().JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	if claims.ExpiresAt.Sub(time.Now().UTC()) < 0 {
		return nil, errors.New("token is expired")
	}
	if claims.Audience[0] != jwtProvider.configuration.Authentication().AudienceField ||
		claims.Issuer != jwtProvider.configuration.Authentication().Issuer {
		return nil, errors.New("invalid token")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}
	userIdParsed := shared.UserId(userId)
	return &userIdParsed, nil
}
