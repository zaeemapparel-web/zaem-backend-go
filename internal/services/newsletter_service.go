package services

import (
	"context"
	"strings"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
)

type NewsletterService struct {
	repo repositories.NewsletterRepository
}

func NewNewsletterService(r repositories.NewsletterRepository) *NewsletterService {
	return &NewsletterService{repo: r}
}

func (s *NewsletterService) Subscribe(ctx context.Context, req dtos.NewsletterSubscribeRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.Internal("Failed to subscribe")
	}

	if existing != nil {
		if existing.IsActive {
			return apperrors.Conflict("Email is already subscribed")
		}
		// Re-activate
		return s.repo.UpdateActive(ctx, existing.ID, true)
	}

	n := &models.Newsletter{
		ID:       cuid.New(),
		Email:    email,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return apperrors.Internal("Failed to subscribe")
	}

	return nil
}

func (s *NewsletterService) Unsubscribe(ctx context.Context, req dtos.NewsletterUnsubscribeRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil || existing == nil {
		return nil // Silent success for security
	}

	return s.repo.UpdateActive(ctx, existing.ID, false)
}