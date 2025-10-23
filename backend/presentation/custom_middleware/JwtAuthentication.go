package customMiddleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/entities"
	"xsedox.com/main/presentation/presentationErrors"
)

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			presentationErrors.WriteJsonFailure(w, "Missing Auth", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		type claimsStruct struct {
			GivenName  string            `json:"given_name" validate:"required"`
			FamilyName string            `json:"family_name" validate:"required"`
			RoomId     *uuid.UUID        `json:"room_id"`
			Roles      entities.UserRole `json:"roles" validate:"required,user_role_validation"`
			jwt.RegisteredClaims
		}

		token, err := jwt.ParseWithClaims(tokenString, &claimsStruct{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.Authentication.JwtKey), nil
		})
		if err != nil {
			presentationErrors.WriteJsonFailure(w, "Invalid Auth", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*claimsStruct)
		if !ok || !token.Valid {
			presentationErrors.WriteJsonFailure(w, "Invalid Auth", http.StatusUnauthorized)
			return
		}
		ctxWithClaims := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctxWithClaims))
	})
}
