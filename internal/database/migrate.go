package database

import (
	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/models"
)

func AutoMigrate() {
	log.Info().Msg("🔄 Running AutoMigrate...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.ReviewHelpful{},
		&models.WishlistItem{},
		&models.Coupon{},
		&models.CouponUsage{},
		&models.LoyaltyTransaction{},
		&models.Notification{},
		&models.Newsletter{},
		&models.EmailLog{},
		&models.AuditLog{},
	)

	if err != nil {
		log.Fatal().Err(err).Msg("❌ AutoMigrate failed")
	}

	log.Info().Msg("✅ All tables migrated successfully")
}