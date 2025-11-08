package authentication

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/user"
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

func (jwtProvider *JwtProvider) GenerateToken(userId user.Id) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   uuid.UUID(userId).String(),
		Audience:  strings.Split(jwtProvider.configuration.Authentication().AudienceField, " "),
		Issuer:    jwtProvider.configuration.Authentication().Issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpirationTime).UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	})
	tokenString, err := token.SignedString([]byte(jwtProvider.configuration.Authentication().JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (jwtProvider *JwtProvider) ValidateTokenAndGetUserId(tokenString string) (*user.Id, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtProvider.configuration.Authentication().JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
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
	userIdParsed := user.Id(userId)
	return &userIdParsed, nil
}
