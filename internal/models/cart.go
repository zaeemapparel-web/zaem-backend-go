package models

import "time"

type Cart struct {
	ID     string `gorm:"primaryKey;type:text" json:"id"`
	UserID string `gorm:"uniqueIndex;type:text;not null" json:"userId"`

	User  User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Items []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	ID        string  `gorm:"primaryKey;type:text" json:"id"`
	CartID    string  `gorm:"type:text;index;not null" json:"cartId"`
	ProductID string  `gorm:"type:text;index;not null" json:"productId"`
	Quantity  int     `gorm:"default:1" json:"quantity"`
	Size      *string `json:"size,omitempty"`
	Color     *string `json:"color,omitempty"`

	Cart    Cart    `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (CartItem) TableName() string {
	return "cart_items"
}