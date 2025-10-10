package domain

import "time"

type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ListingID int       `json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
