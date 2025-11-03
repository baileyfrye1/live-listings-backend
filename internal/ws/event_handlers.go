package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"server/internal/domain"
)

type EventHandler func(event Event, client *WSClient) error

func handleFavoritedListing(event Event, client *WSClient) error {
	var favoriteListingEvent FavoritedListingEvent
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := json.Unmarshal(event.Payload, &favoriteListingEvent); err != nil {
		return fmt.Errorf("Failed to unmarshal event: %v", err)
	}

	message := fmt.Sprintf("New favorite on %s", favoriteListingEvent.Address)

	agentId, err := client.Manager.NotificationService.GetAgentIdByListingId(
		ctx,
		favoriteListingEvent.ListingID,
	)
	if err != nil {
		return fmt.Errorf("Failed to retrieve agent id: %v", err)
	}

	newNotification := &domain.Notification{
		UserID:    agentId,
		ListingID: favoriteListingEvent.ListingID,
		Type:      EventFavoritedListingNotification,
		Message:   message,
	}

	_, err = client.Manager.NotificationService.CreateNotification(ctx, newNotification)
	if err != nil {
		return fmt.Errorf("Failed to persist notification: %w", err)
	}

	notificationPayload := Notification{Message: message}

	data, err := json.Marshal(notificationPayload)
	if err != nil {
		return fmt.Errorf("Failed to marshal notification: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventFavoritedListingNotification,
	}

	for c := range client.Manager.Clients {
		if c.UserId == agentId {
			c.Egress <- outgoingEvent
		}
	}

	return nil
}

func handlePriceDrop(event Event, client *WSClient) error {
	var priceDropEvent PriceDropEvent
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := json.Unmarshal(event.Payload, &priceDropEvent); err != nil {
		return fmt.Errorf("Failed to unmarshal event: %v", err)
	}

	message := fmt.Sprintf(
		"Price Drop: %s was reduced to %s",
		priceDropEvent.Address,
		priceDropEvent.Price,
	)

	return broadcastMessage(
		ctx,
		client,
		priceDropEvent.ListingID,
		message,
		EventPriceDropNotification,
	)
}

func handleStatusChange(event Event, client *WSClient) error {
	var statusChangeEvent StatusChangeEvent
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := json.Unmarshal(event.Payload, &statusChangeEvent); err != nil {
		return fmt.Errorf("Failed to unmarshal event: %v", err)
	}

	message := fmt.Sprintf(
		"Status Change: Status of %s was changed to %s",
		statusChangeEvent.Address,
		statusChangeEvent.Status,
	)

	return broadcastMessage(
		ctx,
		client,
		statusChangeEvent.ListingID,
		message,
		EventStatusChangeNotification,
	)
}

func broadcastMessage(
	ctx context.Context,
	client *WSClient,
	listingId int,
	message string,
	eventType string,
) error {
	userIds, err := client.Manager.NotificationService.GetAllUserIdsByListingId(
		ctx,
		listingId,
	)
	if err != nil {
		return fmt.Errorf("Failed to fetch users for listing: %w", err)
	}

	if len(userIds) == 0 {
		return nil
	}

	for userId := range userIds {
		newNotification := &domain.Notification{
			UserID:    userId,
			ListingID: listingId,
			Type:      eventType,
			Message:   message,
		}

		_, err = client.Manager.NotificationService.CreateNotification(ctx, newNotification)
		if err != nil {
			return fmt.Errorf("Failed to persist notification: %w", err)
		}
	}

	notificationPayload := Notification{Message: message}

	data, err := json.Marshal(notificationPayload)
	if err != nil {
		return fmt.Errorf("Failed to marshal notification: %w", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    eventType,
	}

	for c := range client.Manager.Clients {
		if _, ok := userIds[c.UserId]; ok {
			c.Egress <- outgoingEvent
		}
	}

	return nil
}
