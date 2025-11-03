package repo

import (
	"context"
	"database/sql"

	"server/internal/domain"
)

type INotificationRepo interface {
	GetAllNotificationsByUserId(
		ctx context.Context,
		userId int,
	) ([]*domain.Notification, error)
	CreateNotification(
		ctx context.Context,
		notification *domain.Notification,
	) (*domain.Notification, error)
	ToggleNotificationReadStatus(ctx context.Context, userId int) (*domain.Notification, error)
}

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) GetAllNotificationsByUserId(
	ctx context.Context,
	userId int,
) ([]*domain.Notification, error) {
	query := `
		SELECT * FROM notifications
		WHERE user_id = $1
	`

	var notifications []*domain.Notification

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		notification := new(domain.Notification)

		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.ListingID,
			&notification.Type,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) CreateNotification(
	ctx context.Context,
	notification *domain.Notification,
) (*domain.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, listing_id, type, message, is_read)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	newNotification := *notification

	err := r.db.QueryRowContext(
		ctx,
		query,
		notification.UserID,
		notification.ListingID,
		notification.Type,
		notification.Message,
		notification.IsRead,
	).Scan(&newNotification.ID, &newNotification.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &newNotification, nil
}

func (r *NotificationRepository) ToggleNotificationReadStatus(
	ctx context.Context,
	id int,
) (*domain.Notification, error) {
	query := `
		UPDATE notifications
		SET is_read = NOT is_read
		WHERE id = $1
		RETURNING *
	`

	var updatedNotification domain.Notification

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&updatedNotification.ID,
		&updatedNotification.UserID,
		&updatedNotification.ListingID,
		&updatedNotification.Type,
		&updatedNotification.Message,
		&updatedNotification.IsRead,
		&updatedNotification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedNotification, nil
}
