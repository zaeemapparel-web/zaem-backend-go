package models

import "time"

type Category struct {
	ID          string  `gorm:"primaryKey;type:text" json:"id"`
	Name        string  `gorm:"not null" json:"name"`
	Slug        string  `gorm:"uniqueIndex;not null" json:"slug"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	ParentID    *string `gorm:"type:text;index" json:"parentId,omitempty"`
	IsParent    bool    `gorm:"default:false" json:"isParent"`
	SortOrder   int     `gorm:"default:0" json:"sortOrder"`
	IsActive    bool    `gorm:"default:true" json:"isActive"`

	Parent   *Category  `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
	Products []Product  `gorm:"foreignKey:CategoryID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Category) TableName() string { return "categories" }