package domain

import "time"

type SessionData struct {
	UserID    int
	CreatedAt time.Time
}

type ContextSessionData struct {
	SessionID string
	UserID    int
	Role      string
}
