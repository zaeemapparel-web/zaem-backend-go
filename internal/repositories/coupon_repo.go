package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository() CouponRepository {
	return &couponRepository{db: database.DB}
}

func (r *couponRepository) FindByCode(ctx context.Context, code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&coupon).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *couponRepository) IncrementUsage(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Coupon{}).
		Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).
		Error
}

func (r *couponRepository) CreateUsage(ctx context.Context, usage *models.CouponUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *couponRepository) CountUsageByUser(ctx context.Context, couponID, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CouponUsage{}).
		Where("coupon_id = ? AND user_id = ?", couponID, userID).
		Count(&count).Error
	return count, err
}

func (r *couponRepository) CountUserOrders(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}