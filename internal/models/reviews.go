package models

import "time"

type Review struct {
	ID             string     `gorm:"primaryKey;type:text" json:"id"`
	UserID         string     `gorm:"type:text;index;not null" json:"userId"`
	ProductID      string     `gorm:"type:text;index;not null" json:"productId"`
	Rating         int        `gorm:"not null" json:"rating"`
	Title          *string    `json:"title,omitempty"`
	Comment        string     `gorm:"not null" json:"comment"`
	Images         []string   `gorm:"type:text[];serializer:json" json:"images"`
	IsApproved     bool       `gorm:"default:false;index" json:"isApproved"`
	IsVerified     bool       `gorm:"default:false" json:"isVerified"`
	HelpfulCount   int        `gorm:"default:0" json:"helpfulCount"`
	AdminReply     *string    `json:"adminReply,omitempty"`
	AdminRepliedAt *time.Time `json:"adminRepliedAt,omitempty"`

	User    User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Product Product         `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"-"`
	Helpful []ReviewHelpful `gorm:"foreignKey:ReviewID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Review) TableName() string {
	return "reviews"
}

type ReviewHelpful struct {
	ID       string `gorm:"primaryKey;type:text" json:"id"`
	ReviewID string `gorm:"type:text;index;not null" json:"reviewId"`
	UserID   string `gorm:"type:text;index;not null" json:"userId"`

	Review Review `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"-"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

func (ReviewHelpful) TableName() string {
	return "review_helpful"
}