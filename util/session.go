package util

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"server/internal/domain"
	"server/internal/session"
)

func CreateSession(
	ctx context.Context,
	session session.ISession,
	user *domain.User,
) (string, error) {
	sessionID := generateSessionID()
	sessionData := &domain.SessionData{
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}

	if err := session.SetSession(ctx, sessionID, sessionData, time.Hour*24); err != nil {
		return "", fmt.Errorf("Error creating user session: %v", err)
	}

	return sessionID, nil
}

func generateSessionID() string {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		slog.Error("Failed to generate session id", slog.String("error", err.Error()))
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
