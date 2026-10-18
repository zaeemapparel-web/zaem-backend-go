package models

import "time"

type Product struct {
	ID           string   `gorm:"primaryKey;type:text" json:"id"`
	Name         string   `gorm:"not null;index" json:"name"`
	Slug         string   `gorm:"uniqueIndex;not null" json:"slug"`
	Description  string   `gorm:"not null" json:"description"`
	Price        float64  `gorm:"not null;index" json:"price"`
	ComparePrice *float64 `json:"comparePrice,omitempty"`
	CategoryID   string   `gorm:"type:text;index;not null" json:"categoryId"`

	Images []string `gorm:"type:text[];serializer:json" json:"images"`
	Sizes  []string `gorm:"type:text[];serializer:json" json:"sizes"`
	Colors []string `gorm:"type:text[];serializer:json" json:"colors"`
	Tags   []string `gorm:"type:text[];serializer:json" json:"tags"`

	Stock      int     `gorm:"default:0;index" json:"stock"`
	SKU        *string `gorm:"uniqueIndex" json:"sku,omitempty"`
	IsFeatured bool    `gorm:"default:false;index" json:"isFeatured"`
	IsActive   bool    `gorm:"default:true;index" json:"isActive"`

	ViewCount int `gorm:"default:0" json:"viewCount"`
	SoldCount int `gorm:"default:0" json:"soldCount"`

	Category  Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CartItems []CartItem     `gorm:"foreignKey:ProductID" json:"-"`
	Reviews   []Review       `gorm:"foreignKey:ProductID" json:"-"`
	Wishlist  []WishlistItem `gorm:"foreignKey:ProductID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Product) TableName() string { return "products" }