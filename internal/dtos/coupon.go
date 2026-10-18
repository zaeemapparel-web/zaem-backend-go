package dtos

import "time"

// ==================== REQUESTS ====================

type ValidateCouponRequest struct {
	Code     string  `json:"code" binding:"required"`
	Subtotal float64 `json:"subtotal" binding:"required,gt=0"`
}

type CreateCouponRequest struct {
	Code                 string     `json:"code" binding:"required,min=3,max=50"`
	Description          *string    `json:"description"`
	Type                 string     `json:"type" binding:"required,oneof=PERCENTAGE FIXED FREE_SHIPPING BUY_X_GET_Y"`
	Value                float64    `json:"value" binding:"required,gt=0"`
	MinOrderAmount       float64    `json:"minOrderAmount"`
	MaxDiscount          *float64   `json:"maxDiscount"`
	MaxUses              *int       `json:"maxUses"`
	MaxUsesPerUser       int        `json:"maxUsesPerUser"`
	BuyQuantity          *int       `json:"buyQuantity"`
	GetQuantity          *int       `json:"getQuantity"`
	IsActive             *bool      `json:"isActive"`
	IsFirstOrder         bool       `json:"isFirstOrder"`
	ExpiresAt            *time.Time `json:"expiresAt"`
	ApplicableCategories []string   `json:"applicableCategories"`
	ApplicableProducts   []string   `json:"applicableProducts"`
}

// ==================== RESPONSES ====================

type CouponValidationResponse struct {
	Valid          bool    `json:"valid"`
	Code           string  `json:"code"`
	Type           string  `json:"type"`
	DiscountAmount float64 `json:"discountAmount"`
	Message        string  `json:"message"`
	FreeShipping   bool    `json:"freeShipping"`
}

type CouponResponse struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Description *string  `json:"description,omitempty"`
	Type        string   `json:"type"`
	Value       float64  `json:"value"`
	MinOrderAmount float64 `json:"minOrderAmount"`
	IsActive    bool     `json:"isActive"`
	UsedCount   int      `json:"usedCount"`
	ExpiresAt   *string  `json:"expiresAt,omitempty"`
}