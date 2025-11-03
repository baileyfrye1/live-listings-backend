package repo

import (
	"context"

	"server/internal/domain"
)

type NotificationRepoMock struct {
	GetAllNotificationsByUserIdFunc func(ctx context.Context, userId int) ([]*domain.Notification, error)
	CreateNotificationFunc          func(ctx context.Context, notification *domain.Notification) (*domain.Notification, error)
	MarkNotificationAsReadFunc      func(ctx context.Context, userId int) (*domain.Notification, error)
}

func (n *NotificationRepoMock) GetAllNotificationsByUserId(
	ctx context.Context,
	userId int,
) ([]*domain.Notification, error) {
	return n.GetAllNotificationsByUserIdFunc(ctx, userId)
}

func (n *NotificationRepoMock) CreateNotification(
	ctx context.Context,
	notification *domain.Notification,
) (*domain.Notification, error) {
	return n.CreateNotificationFunc(ctx, notification)
}

func (n *NotificationRepoMock) MarkNotificationAsRead(
	ctx context.Context,
	userId int,
) (*domain.Notification, error) {
	return n.MarkNotificationAsReadFunc(ctx, userId)
}
