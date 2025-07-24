package domain

import "time"

type SessionData struct {
	UserID    int
	CreatedAt time.Time
}
