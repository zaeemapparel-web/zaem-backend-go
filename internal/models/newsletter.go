package models

import "time"

type Newsletter struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	IsActive  bool      `gorm:"default:true" json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Newsletter) TableName() string { return "newsletters" }