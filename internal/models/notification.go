package models

import "time"

type Notification struct {
	ID       string           `gorm:"primaryKey;type:text" json:"id"`
	UserID   string           `gorm:"type:text;index;not null" json:"userId"`
	Type     NotificationType `gorm:"type:text;not null" json:"type"`
	Title    string           `gorm:"not null" json:"title"`
	Message  string           `gorm:"not null" json:"message"`
	Link     *string          `json:"link,omitempty"`
	ImageURL *string          `json:"imageUrl,omitempty"`
	IsRead   bool             `gorm:"default:false;index" json:"isRead"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

func (Notification) TableName() string { return "notifications" }