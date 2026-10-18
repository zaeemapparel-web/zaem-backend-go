package models

import "time"

type Coupon struct {
	ID                   string     `gorm:"primaryKey;type:text" json:"id"`
	Code                 string     `gorm:"uniqueIndex;not null" json:"code"`
	Description          *string    `json:"description,omitempty"`
	Type                 CouponType `gorm:"type:text;not null" json:"type"`
	Value                float64    `gorm:"not null" json:"value"`
	MinOrderAmount       float64    `gorm:"default:0" json:"minOrderAmount"`
	MaxDiscount          *float64   `json:"maxDiscount,omitempty"`
	MaxUses              *int       `json:"maxUses,omitempty"`
	UsedCount            int        `gorm:"default:0" json:"usedCount"`
	MaxUsesPerUser       int        `gorm:"default:1" json:"maxUsesPerUser"`
	BuyQuantity          *int       `json:"buyQuantity,omitempty"`
	GetQuantity          *int       `json:"getQuantity,omitempty"`
	IsActive             bool       `gorm:"default:true;index" json:"isActive"`
	IsFirstOrder         bool       `gorm:"default:false" json:"isFirstOrder"`
	StartsAt             time.Time  `gorm:"default:now()" json:"startsAt"`
	ExpiresAt            *time.Time `json:"expiresAt,omitempty"`
	ApplicableCategories []string   `gorm:"type:text[];serializer:json" json:"applicableCategories"`
	ApplicableProducts   []string   `gorm:"type:text[];serializer:json" json:"applicableProducts"`

	Orders []Order       `gorm:"foreignKey:CouponID" json:"-"`
	Usage  []CouponUsage `gorm:"foreignKey:CouponID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Coupon) TableName() string { return "coupons" }

type CouponUsage struct {
	ID       string  `gorm:"primaryKey;type:text" json:"id"`
	CouponID string  `gorm:"type:text;index;not null" json:"couponId"`
	UserID   string  `gorm:"type:text;index;not null" json:"userId"`
	OrderID  *string `gorm:"type:text" json:"orderId,omitempty"`

	Coupon Coupon `gorm:"foreignKey:CouponID;constraint:OnDelete:CASCADE" json:"-"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

func (CouponUsage) TableName() string { return "coupon_usage" }