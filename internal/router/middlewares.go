package router

import (
	"context"
	"net/http"

	"github.com/Vykiy/house-service/internal/models"
	"github.com/google/uuid"
)

const userIDCtxKey = "user_id"

type Middleware struct {
	jwtIssuer *JWTIssuer
}

func NewMiddleware(jwtIssuer *JWTIssuer) *Middleware {
	return &Middleware{jwtIssuer: jwtIssuer}
}

func (m *Middleware) UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, _, err := m.parseUserFromHeader(r.Header)
		if err != nil {
			http.Error(w, "неверный токен", http.StatusUnauthorized)
			return
		}

		if userType != models.UserTypeUser && userType != models.UserTypeModerator {
			http.Error(w, "недостаточно прав", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ModeratorAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, userID, err := m.parseUserFromHeader(r.Header)
		if err != nil {
			http.Error(w, "неверный токен", http.StatusUnauthorized)
			return
		}

		if userType != models.UserTypeModerator {
			http.Error(w, "недостаточно прав", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), userIDCtxKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) parseUserFromHeader(header http.Header) (models.UserType, uuid.UUID, error) {
	jwt := header.Get("Authorization")
	if jwt == "" {
		return models.UserTypeUnknown, uuid.Nil, nil
	}

	userType, userID, err := m.jwtIssuer.ParseToken(jwt)
	if err != nil {
		return models.UserTypeUnknown, uuid.Nil, err
	}

	return userType, userID, nil
}
