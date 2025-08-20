package dto

type CreateListingRequest struct {
	Address     string  `json:"address"`
	Price       int     `json:"price"`
	Beds        int     `json:"beds"`
	Baths       int     `json:"baths"`
	SqFt        int     `json:"sq_ft"`
	Description *string `json:"description"`
	AgentID     *int    `json:"agent_id"`
}

type UpdateListingRequest struct {
	Address     *string `json:"address"`
	Price       *int    `json:"price"`
	Beds        *int    `json:"beds"`
	Baths       *int    `json:"baths"`
	SqFt        *int    `json:"sq_ft"`
	Description *string `json:"description"`
	AgentID     *int    `json:"agent_id"`
}
