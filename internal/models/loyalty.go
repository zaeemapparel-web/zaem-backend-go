package models

import "time"

type LoyaltyTransaction struct {
	ID          string        `gorm:"primaryKey;type:text" json:"id"`
	UserID      string        `gorm:"type:text;index;not null" json:"userId"`
	Points      int           `gorm:"not null" json:"points"`
	Action      LoyaltyAction `gorm:"type:text;not null" json:"action"`
	Description string        `gorm:"not null" json:"description"`
	OrderID     *string       `json:"orderId,omitempty"`
	ReviewID    *string       `json:"reviewId,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"index" json:"createdAt"`
}

func (LoyaltyTransaction) TableName() string { return "loyalty_transactions" }