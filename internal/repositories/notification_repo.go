package repositories

import (
	"context"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
)

type notificationRepository struct {
	db interface{}
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]models.Notification, int64, error) {
	query := database.DB.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	var items []models.Notification
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&items).Error

	return items, total, err
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id, userID string) error {
	return database.DB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).
		Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return database.DB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).
		Error
}

func (r *notificationRepository) Delete(ctx context.Context, id, userID string) error {
	return database.DB.WithContext(ctx).
		Delete(&models.Notification{}, "id = ? AND user_id = ?", id, userID).
		Error
}

func (r *notificationRepository) Create(ctx context.Context, n *models.Notification) error {
	return database.DB.WithContext(ctx).Create(n).Error
}