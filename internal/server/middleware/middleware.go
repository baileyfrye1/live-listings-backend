package middleware

import (
	"context"
	"net/http"
	"net/url"

	"server/internal/repo"
	"server/internal/session"
)

type contextKey string

const (
	UserContextKey contextKey = "userID"
	RoleContextKey contextKey = "role"
)

func Authenticate(
	session *session.Session,
	userRepo *repo.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie("session")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionID, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionData, err := session.GetSession(ctx, sessionID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, err = userRepo.GetUserById(ctx, sessionData.UserID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userCtx := context.WithValue(ctx, UserContextKey, sessionData.UserID)
			next.ServeHTTP(w, r.WithContext(userCtx))
		})
	}
}

func Authorize(
	session *session.Session,
	userRepo *repo.UserRepository,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cookie, err := r.Cookie("session")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionID, err := url.QueryUnescape(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionData, err := session.GetSession(ctx, sessionID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetAuthorizedUser(ctx, sessionData.UserID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userCtx := context.WithValue(ctx, UserContextKey, sessionData.UserID)
			roleCtx := context.WithValue(userCtx, RoleContextKey, user.Role)
			next.ServeHTTP(w, r.WithContext(roleCtx))
		})
	}
}
