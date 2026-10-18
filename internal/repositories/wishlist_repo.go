package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository() WishlistRepository {
	return &wishlistRepository{db: database.DB}
}

func (r *wishlistRepository) FindByUserID(ctx context.Context, userID string) ([]models.WishlistItem, error) {
	var items []models.WishlistItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *wishlistRepository) FindByUserAndProduct(ctx context.Context, userID, productID string) (*models.WishlistItem, error) {
	var item models.WishlistItem
	err := r.db.WithContext(ctx).
		First(&item, "user_id = ? AND product_id = ?", userID, productID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *wishlistRepository) Create(ctx context.Context, item *models.WishlistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *wishlistRepository) Delete(ctx context.Context, userID, productID string) error {
	return r.db.WithContext(ctx).
		Delete(&models.WishlistItem{}, "user_id = ? AND product_id = ?", userID, productID).Error
}

func (r *wishlistRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.WishlistItem{}, "id = ?", id).Error
}

func (r *wishlistRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.WishlistItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}