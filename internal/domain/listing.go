package domain

import "time"

type Listing struct {
	ID          int           `json:"id"`
	Address     string        `json:"address"`
	Price       int           `json:"price"`
	Beds        int           `json:"beds"`
	Baths       int           `json:"baths"`
	SqFt        int           `json:"sq_ft"`
	Description *string       `json:"description"`
	AgentID     int           `json:"agent_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Agent       *ListingAgent `json:"agent"`
}

type ListingAgent struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
