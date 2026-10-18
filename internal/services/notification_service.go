package services

import (
	"context"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/repositories"
)

type NotificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(r repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: r}
}

func (s *NotificationService) GetMyNotifications(ctx context.Context, userID string, page, limit int) ([]dtos.NotificationResponse, int64, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 20 }

	items, total, err := s.repo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, 0, apperrors.Internal("Failed to fetch notifications")
	}

	unread, _ := s.repo.CountUnread(ctx, userID)

	result := make([]dtos.NotificationResponse, 0, len(items))
	for _, n := range items {
		result = append(result, dtos.NotificationResponse{
			ID:        n.ID,
			Type:      string(n.Type),
			Title:     n.Title,
			Message:   n.Message,
			Link:      n.Link,
			ImageURL:  n.ImageURL,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, total, unread, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, userID, id string) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}