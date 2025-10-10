package middleware

import (
	"context"
	"net/http"
	"net/url"

	"server/internal/domain"
)

type SessionProvider interface {
	GetSession(ctx context.Context, sessionId string) (*domain.SessionData, error)
}

type UserProvider interface {
	GetUserById(ctx context.Context, userId int) (*domain.User, error)
}

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func Authenticate(
	session SessionProvider,
	userRepo UserProvider,
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

			user, err := userRepo.GetUserById(ctx, sessionData.UserID)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			contextSessionData := &domain.ContextSessionData{
				SessionID: sessionID,
				UserID:    sessionData.UserID,
				Role:      user.Role,
			}

			userCtx := context.WithValue(ctx, UserContextKey, contextSessionData)
			next.ServeHTTP(w, r.WithContext(userCtx))
		})
	}
}

func Authorize() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sess := ctx.Value(UserContextKey).(*domain.ContextSessionData)

			if sess.Role == "user" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
