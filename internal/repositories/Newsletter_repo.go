package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type newsletterRepository struct{}

func NewNewsletterRepository() NewsletterRepository {
	return &newsletterRepository{}
}

func (r *newsletterRepository) FindByEmail(ctx context.Context, email string) (*models.Newsletter, error) {
	var n models.Newsletter
	err := database.DB.WithContext(ctx).First(&n, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *newsletterRepository) Create(ctx context.Context, n *models.Newsletter) error {
	return database.DB.WithContext(ctx).Create(n).Error
}

func (r *newsletterRepository) UpdateActive(ctx context.Context, id string, isActive bool) error {
	return database.DB.WithContext(ctx).
		Model(&models.Newsletter{}).
		Where("id = ?", id).
		Update("is_active", isActive).
		Error
}