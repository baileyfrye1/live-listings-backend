package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/domain"
)

type sessionMock struct {
	GetSessionFunc func(ctx context.Context, sessionId string) (*domain.SessionData, error)
}

func (m *sessionMock) GetSession(
	ctx context.Context,
	sessionId string,
) (*domain.SessionData, error) {
	return m.GetSessionFunc(ctx, sessionId)
}

type userRepoMock struct {
	GetByUserIdFunc func(ctx context.Context, id int) (*domain.User, error)
}

func (u *userRepoMock) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	return u.GetByUserIdFunc(ctx, id)
}

func TestAuthentication(t *testing.T) {
	cases := []struct {
		description    string
		cookie         *http.Cookie
		mockSession    func() *sessionMock
		mockRepo       func() *userRepoMock
		wantStatusCode int
		wantNextCalled bool
	}{
		{
			description: "Valid cookie and session returns success",
			cookie:      &http.Cookie{Name: "session", Value: "abc123"},
			mockSession: func() *sessionMock {
				return &sessionMock{
					GetSessionFunc: func(ctx context.Context, sessionId string) (*domain.SessionData, error) {
						return &domain.SessionData{UserID: 123, CreatedAt: time.Now()}, nil
					},
				}
			},
			mockRepo: func() *userRepoMock {
				return &userRepoMock{
					GetByUserIdFunc: func(ctx context.Context, id int) (*domain.User, error) {
						return &domain.User{ID: 123, Role: "agent"}, nil
					},
				}
			},
			wantStatusCode: http.StatusOK,
			wantNextCalled: true,
		},
		{
			description: "Missing cookie returns unauthorized",
			cookie:      nil,
			mockSession: func() *sessionMock {
				return &sessionMock{
					GetSessionFunc: func(ctx context.Context, sessionId string) (*domain.SessionData, error) {
						return &domain.SessionData{UserID: 123, CreatedAt: time.Now()}, nil
					},
				}
			},
			mockRepo: func() *userRepoMock {
				return &userRepoMock{
					GetByUserIdFunc: func(ctx context.Context, id int) (*domain.User, error) {
						return &domain.User{ID: 123, Role: "agent"}, nil
					},
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			description: "Session lookup fails returns unauthorized",
			cookie:      &http.Cookie{Name: "session", Value: "abc123"},
			mockSession: func() *sessionMock {
				return &sessionMock{
					GetSessionFunc: func(ctx context.Context, sessionId string) (*domain.SessionData, error) {
						return nil, errors.New("Session not found")
					},
				}
			},
			mockRepo: func() *userRepoMock {
				return &userRepoMock{
					GetByUserIdFunc: func(ctx context.Context, id int) (*domain.User, error) {
						return &domain.User{ID: 123, Role: "agent"}, nil
					},
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			description: "User lookup fails returns unauthorized",
			cookie:      &http.Cookie{Name: "session", Value: "abc123"},
			mockSession: func() *sessionMock {
				return &sessionMock{
					GetSessionFunc: func(ctx context.Context, sessionId string) (*domain.SessionData, error) {
						return &domain.SessionData{UserID: 123, CreatedAt: time.Now()}, nil
					},
				}
			},
			mockRepo: func() *userRepoMock {
				return &userRepoMock{
					GetByUserIdFunc: func(ctx context.Context, id int) (*domain.User, error) {
						return nil, errors.New("User not found")
					},
				}
			},
			wantStatusCode: http.StatusUnauthorized,
			wantNextCalled: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			mw := Authenticate(tt.mockSession(), tt.mockRepo())
			mw(next).ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatusCode)
			}

			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNextCalled)
			}
		})
	}
}

func TestAuthorization(t *testing.T) {
	cases := []struct {
		description    string
		role           string
		wantStatusCode int
		wantNextCalled bool
	}{
		{
			description:    "Agent role allowed returns success",
			role:           "agent",
			wantStatusCode: http.StatusOK,
			wantNextCalled: true,
		},
		{
			description:    "Admin role allowed returns success",
			role:           "admin",
			wantStatusCode: http.StatusOK,
			wantNextCalled: true,
		},
		{
			description:    "User role not allowed returns forbidden",
			role:           "user",
			wantStatusCode: http.StatusForbidden,
			wantNextCalled: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			ctx := context.WithValue(
				req.Context(),
				UserContextKey,
				&domain.ContextSessionData{
					UserID: 123,
					Role:   tt.role,
				},
			)
			req = req.WithContext(ctx)

			mw := Authorize()
			mw(next).ServeHTTP(rr, req)

			if rr.Code != tt.wantStatusCode {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatusCode)
			}
			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNextCalled)
			}
		})
	}
}
