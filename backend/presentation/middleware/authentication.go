package customMiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"xsedox.com/main/application/response"
	"xsedox.com/main/config"
	"xsedox.com/main/domain/entities"
	"xsedox.com/main/presentation/errors"
	"xsedox.com/main/validation"
)

type Claims struct {
	GivenName  string            `json:"given_name" validate:"required"`
	FamilyName string            `json:"family_name" validate:"required"`
	RoomId     *uuid.UUID        `json:"room_id"`
	Roles      entities.UserRole `json:"roles" validate:"required,user_role_validation"`
	jwt.RegisteredClaims
}

func Authentication(next http.Handler, conf *config.Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(response.Failure("Authorization header required")); err != nil {
				http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
			}
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			w.WriteHeader(http.StatusUnauthorized)

			if err := json.NewEncoder(w).Encode(response.Failure("Could not find bearer token in Authorization header")); err != nil {
				http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
			}
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(conf.Authentication.JwtSecret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				if innerr := json.NewEncoder(w).Encode(response.Failure("Invalid token signature")); innerr != nil {
					http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
					return
				}
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			if inerr := json.NewEncoder(w).Encode(response.Failure("Invalid token")); inerr != nil {
				http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
			}
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(response.Failure("Invalid token")); err != nil {
				http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
			}
			return
		}
		if inerr := validation.Validate.Struct(claims); inerr != nil {
			w.WriteHeader(http.StatusBadRequest)
			// Provide a more specific error message
			if err := json.NewEncoder(w).Encode(response.Failure("Token claims validation failed")); err != nil {
				http.Error(w, errors.EncodingErrorMessage, http.StatusInternalServerError)
			}
			return
		}

		// Add claims to the context
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
