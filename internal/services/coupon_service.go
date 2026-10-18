package services

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"

	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
)

type CouponService struct {
	repo repositories.CouponRepository
}

func NewCouponService(r repositories.CouponRepository) *CouponService {
	return &CouponService{repo: r}
}

// ValidateForCheckout — used during order creation
func (s *CouponService) ValidateForCheckout(ctx context.Context, code string, userID string, subtotal float64) (*models.Coupon, float64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	coupon, err := s.repo.FindByCode(ctx, code)
	if err != nil || coupon == nil {
		return nil, 0, apperrors.BadRequest("Invalid coupon code")
	}

	// Check active
	if !coupon.IsActive {
		return nil, 0, apperrors.BadRequest("Coupon is not active")
	}

	// Check expiry
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now().UTC()) {
		return nil, 0, apperrors.BadRequest("Coupon has expired")
	}

	// Check max uses
	if coupon.MaxUses != nil && coupon.UsedCount >= *coupon.MaxUses {
		return nil, 0, apperrors.BadRequest("Coupon usage limit reached")
	}

	// Check min order
	if subtotal < coupon.MinOrderAmount {
		return nil, 0, apperrors.BadRequest("Order amount too low for this coupon")
	}

	// Check per-user usage
	usage, _ := s.repo.CountUsageByUser(ctx, coupon.ID, userID)
	if int(usage) >= coupon.MaxUsesPerUser {
		return nil, 0, apperrors.BadRequest("You have already used this coupon")
	}

	// Check first order
	if coupon.IsFirstOrder {
		orders, _ := s.repo.CountUserOrders(ctx, userID)
		if orders > 0 {
			return nil, 0, apperrors.BadRequest("This coupon is only for first order")
		}
	}

	// Calculate discount
	discount := s.calculateDiscount(coupon, subtotal)

	return coupon, discount, nil
}

// CalculateDiscount for public preview
func (s *CouponService) CalculateDiscount(coupon *models.Coupon, subtotal float64) float64 {
	return s.calculateDiscount(coupon, subtotal)
}

func (s *CouponService) calculateDiscount(coupon *models.Coupon, subtotal float64) float64 {
	var discount float64

	switch coupon.Type {
	case models.CouponPercentage:
		discount = subtotal * (coupon.Value / 100.0)
		if coupon.MaxDiscount != nil && discount > *coupon.MaxDiscount {
			discount = *coupon.MaxDiscount
		}
	case models.CouponFixed:
		discount = coupon.Value
		if discount > subtotal {
			discount = subtotal
		}
	case models.CouponFreeShipping:
		// Handled separately (shipping = 0)
		discount = 0
	case models.CouponBuyXGetY:
		// Complex logic — simplified here
		discount = 0
	}

	return discount
}