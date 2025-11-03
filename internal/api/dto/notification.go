package dto

type CreateNotificationRequest struct {
	ListingID int    `json:"listing_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}
