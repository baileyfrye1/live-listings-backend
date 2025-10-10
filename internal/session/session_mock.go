package session

import (
	"context"
	"time"

	"server/internal/domain"
)

type SessionMock struct {
	GetSessionFunc func(ctx context.Context, sessionID string) (*domain.SessionData, error)
	SetSessionFunc func(
		ctx context.Context,
		sessionID string,
		sessionData *domain.SessionData,
		ttl time.Duration,
	) error
	DeleteSessionFunc func(ctx context.Context, sessionID string) error
}

func (s *SessionMock) GetSession(
	ctx context.Context,
	sessionID string,
) (*domain.SessionData, error) {
	return s.GetSessionFunc(ctx, sessionID)
}

func (s *SessionMock) SetSession(
	ctx context.Context,
	sessionID string,
	sessionData *domain.SessionData,
	ttl time.Duration,
) error {
	return s.SetSessionFunc(ctx, sessionID, sessionData, ttl)
}

func (s *SessionMock) DeleteSession(ctx context.Context, sessionID string) error {
	return s.DeleteSessionFunc(ctx, sessionID)
}
