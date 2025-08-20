package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"server/internal/domain"
)

type Session struct {
	rc *redis.Client
}

func NewSession(rc *redis.Client) *Session {
	return &Session{
		rc: rc,
	}
}

func (s *Session) GetSession(ctx context.Context, sessionID string) (*domain.SessionData, error) {
	val, err := s.rc.Get(ctx, sessionID).Result()
	if err != nil {
		return nil, err
	}

	var sessionData domain.SessionData

	if err = json.Unmarshal([]byte(val), &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (s *Session) SetSession(
	ctx context.Context,
	sessionID string,
	sessionData *domain.SessionData,
	ttl time.Duration,
) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return s.rc.Set(ctx, sessionID, data, ttl).Err()
}

func (s *Session) DeleteSession(ctx context.Context, sessionID string) error {
	count, err := s.rc.Del(ctx, sessionID).Result()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Session not found")
	}

	return nil
}
