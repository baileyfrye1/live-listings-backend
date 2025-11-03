package ws

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ListingEvent struct {
	ListingID int    `json:"listing_id"`
	Address   string `json:"address"`
}

type FavoritedListingEvent struct {
	ListingEvent
}

type PriceDropEvent struct {
	Price string `json:"price"`
	ListingEvent
}

type StatusChangeEvent struct {
	Status string `json:"status"`
	ListingEvent
}

type Notification struct {
	Message string `json:"message"`
}

const (
	EventFavoritedListingNotification = "favorited_listing_notification"
	EventPriceDropNotification        = "price_drop_notification"
	EventStatusChangeNotification     = "status_changed_notification"
)
