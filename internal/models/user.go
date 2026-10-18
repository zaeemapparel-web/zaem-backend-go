package models

import "time"

type User struct {
	ID       string  `gorm:"primaryKey;type:text" json:"id"`
	Name     string  `gorm:"not null" json:"name"`
	Email    string  `gorm:"uniqueIndex;not null" json:"email"`
	Phone    *string `gorm:"uniqueIndex" json:"phone,omitempty"`
	Password string  `gorm:"not null" json:"-"`
	Role     Role    `gorm:"type:text;default:USER" json:"role"`
	IsActive bool    `gorm:"default:true" json:"isActive"`
	Avatar   *string `json:"avatar,omitempty"`

	DateOfBirth *time.Time `json:"dateOfBirth,omitempty"`

	// Email Verification
	IsEmailVerified bool       `gorm:"default:false" json:"isEmailVerified"`
	EmailVerifyToken *string    `gorm:"index" json:"-"`
	EmailVerifiedAt  *time.Time `json:"emailVerifiedAt,omitempty"`

	// Password Reset
	ResetToken       *string    `gorm:"index" json:"-"`
	ResetTokenExpiry *time.Time `json:"-"`

	// Refresh Token
	RefreshToken       *string    `gorm:"index" json:"-"`
	RefreshTokenExpiry *time.Time `json:"-"`

	// Loyalty
	LoyaltyPoints int         `gorm:"default:0" json:"loyaltyPoints"`
	TotalSpent    float64     `gorm:"default:0" json:"totalSpent"`
	Tier          LoyaltyTier `gorm:"type:text;default:BRONZE" json:"tier"`
	ReferralCode  *string     `gorm:"uniqueIndex" json:"referralCode,omitempty"`
	ReferredBy    *string     `json:"referredBy,omitempty"`

	// Last login
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	LastLoginIP *string    `json:"-"`

	// Relations
	Addresses      []Address            `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
	Orders         []Order              `gorm:"foreignKey:UserID" json:"-"`
	Reviews        []Review             `gorm:"foreignKey:UserID" json:"-"`
	Cart           *Cart                `gorm:"foreignKey:UserID" json:"-"`
	Wishlist       []WishlistItem       `gorm:"foreignKey:UserID" json:"-"`
	LoyaltyHistory []LoyaltyTransaction `gorm:"foreignKey:UserID" json:"-"`
	Notifications  []Notification       `gorm:"foreignKey:UserID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string { return "users" }