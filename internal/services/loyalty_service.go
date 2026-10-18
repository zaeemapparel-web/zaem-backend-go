package services

import (
	"context"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/pkg/cuid"
)

type LoyaltyService struct{}

func NewLoyaltyService() *LoyaltyService {
	return &LoyaltyService{}
}

// Points earned = 1 point per Rs. 100 spent
func CalculatePointsEarned(amount float64) int {
	return int(amount / 100)
}

// Add points and update tier
func (s *LoyaltyService) AwardPoints(ctx context.Context, userID string, points int, action models.LoyaltyAction, description string, orderID *string) error {
	if points <= 0 {
		return nil
	}

	// Create transaction
	txn := &models.LoyaltyTransaction{
		ID:          cuid.New(),
		UserID:      userID,
		Points:      points,
		Action:      action,
		Description: description,
		OrderID:     orderID,
	}

	if err := database.DB.WithContext(ctx).Create(txn).Error; err != nil {
		return err
	}

	// Update user points
	if err := database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("loyalty_points", database.DB.Raw("loyalty_points + ?", points)).
		Error; err != nil {
		return err
	}

	// Auto-update tier
	s.UpdateTier(ctx, userID)

	return nil
}

// Redeem points
func (s *LoyaltyService) RedeemPoints(ctx context.Context, userID string, points int, orderID string) error {
	if points <= 0 {
		return nil
	}

	txn := &models.LoyaltyTransaction{
		ID:          cuid.New(),
		UserID:      userID,
		Points:      -points,
		Action:      models.LoyaltyRedeemed,
		Description: "Points redeemed at checkout",
		OrderID:     &orderID,
	}

	if err := database.DB.WithContext(ctx).Create(txn).Error; err != nil {
		return err
	}

	return database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		UpdateColumn("loyalty_points", database.DB.Raw("loyalty_points - ?", points)).
		Error
}

// UpdateTier based on total spent
func (s *LoyaltyService) UpdateTier(ctx context.Context, userID string) {
	var user models.User
	if err := database.DB.WithContext(ctx).Select("id", "total_spent", "tier").First(&user, "id = ?", userID).Error; err != nil {
		return
	}

	newTier := models.TierBronze
	switch {
	case user.TotalSpent >= 500000: // 5 lakh+
		newTier = models.TierPlatinum
	case user.TotalSpent >= 200000: // 2 lakh+
		newTier = models.TierGold
	case user.TotalSpent >= 50000: // 50k+
		newTier = models.TierSilver
	}

	if user.Tier != newTier {
		database.DB.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", userID).
			Update("tier", newTier)
	}
}