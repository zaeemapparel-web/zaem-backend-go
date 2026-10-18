package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository() ReviewRepository {
	return &reviewRepository{db: database.DB}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) FindByID(ctx context.Context, id string) (*models.Review, error) {
	var rev models.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&rev, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rev, err
}

func (r *reviewRepository) FindByUserAndProduct(ctx context.Context, userID, productID string) (*models.Review, error) {
	var rev models.Review
	err := r.db.WithContext(ctx).
		First(&rev, "user_id = ? AND product_id = ?", userID, productID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rev, err
}

func (r *reviewRepository) FindByProductID(ctx context.Context, productID string, page, limit int, onlyApproved bool) ([]models.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Review{}).Where("product_id = ?", productID)

	if onlyApproved {
		query = query.Where("is_approved = ?", true)
	}

	var total int64
	query.Count(&total)

	var reviews []models.Review
	err := query.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *reviewRepository) FindAll(ctx context.Context, page, limit int, isApproved *bool) ([]models.Review, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Review{})

	if isApproved != nil {
		query = query.Where("is_approved = ?", *isApproved)
	}

	var total int64
	query.Count(&total)

	var reviews []models.Review
	err := query.
		Preload("User").
		Preload("Product").
		Order("created_at DESC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}

func (r *reviewRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id).Error
}

// Helpful
func (r *reviewRepository) FindHelpful(ctx context.Context, reviewID, userID string) (*models.ReviewHelpful, error) {
	var h models.ReviewHelpful
	err := r.db.WithContext(ctx).
		First(&h, "review_id = ? AND user_id = ?", reviewID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &h, err
}

func (r *reviewRepository) CreateHelpful(ctx context.Context, h *models.ReviewHelpful) error {
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *reviewRepository) DeleteHelpful(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.ReviewHelpful{}, "id = ?", id).Error
}

func (r *reviewRepository) IncrementHelpfulCount(ctx context.Context, reviewID string, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("id = ?", reviewID).
		UpdateColumn("helpful_count", gorm.Expr("helpful_count + ?", delta)).
		Error
}

// Stats
func (r *reviewRepository) GetProductRatingStats(ctx context.Context, productID string) (float64, int64, error) {
	var result struct {
		Avg   float64
		Count int64
	}
	err := r.db.WithContext(ctx).
		Model(&models.Review{}).
		Where("product_id = ? AND is_approved = ?", productID, true).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Scan(&result).Error
	return result.Avg, result.Count, err
}

// Verify helper
func (r *reviewRepository) HasPurchased(ctx context.Context, userID, productID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ?", userID, productID, models.OrderDelivered).
		Count(&count).Error
	return count > 0, err
}