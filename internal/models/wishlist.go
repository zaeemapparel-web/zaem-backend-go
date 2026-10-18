package models

import "time"

type WishlistItem struct {
	ID            string     `gorm:"primaryKey;type:text" json:"id"`
	UserID        string     `gorm:"type:text;index;not null" json:"userId"`
	ProductID     string     `gorm:"type:text;index;not null" json:"productId"`
	PriceAtAdd    *float64   `json:"priceAtAdd,omitempty"`
	NotifyOnDrop  bool       `gorm:"default:true" json:"notifyOnDrop"`
	NotifyOnStock bool       `gorm:"default:true" json:"notifyOnStock"`
	NotifiedAt    *time.Time `json:"notifiedAt,omitempty"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"`

	CreatedAt time.Time `json:"createdAt"`
}

func (WishlistItem) TableName() string {
	return "wishlist_items"
}