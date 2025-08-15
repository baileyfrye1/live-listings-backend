package domain

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash *string   `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Role         string    `json:"role"`
}
