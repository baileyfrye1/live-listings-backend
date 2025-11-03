package service

import (
	"context"

	"server/internal/domain"
	"server/internal/repo"
)

type NotificationService struct {
	notificationRepo repo.INotificationRepo
	favoriteRepo     repo.IFavoriteRepo
	listingRepo      repo.IListingRepo
}

func NewNotificationService(
	notificationRepo repo.INotificationRepo,
	favoriteRepo repo.IFavoriteRepo,
	listingRepo repo.IListingRepo,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		favoriteRepo:     favoriteRepo,
		listingRepo:      listingRepo,
	}
}

func (s *NotificationService) GetAllNotificationsByUserId(
	ctx context.Context,
	userID int,
) ([]*domain.Notification, error) {
	notifications, err := s.notificationRepo.GetAllNotificationsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	if notifications == nil {
		notifications = []*domain.Notification{}
	}

	return notifications, nil
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	notification *domain.Notification,
) (*domain.Notification, error) {
	return s.notificationRepo.CreateNotification(ctx, notification)
}

func (s *NotificationService) ToggleNotificationReadStatus(
	ctx context.Context,
	id int,
) (*domain.Notification, error) {
	return s.notificationRepo.ToggleNotificationReadStatus(ctx, id)
}

func (s *NotificationService) GetAllUserIdsByListingId(
	ctx context.Context,
	listingId int,
) (map[int]bool, error) {
	return s.favoriteRepo.GetAllUserIdsByListingId(ctx, listingId)
}

func (s *NotificationService) GetAgentIdByListingId(
	ctx context.Context,
	listingId int,
) (int, error) {
	return s.listingRepo.GetAgentIdByListingId(ctx, listingId)
}
