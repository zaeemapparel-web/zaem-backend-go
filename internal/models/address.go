package models

import "time"

type Address struct {
	ID         string  `gorm:"primaryKey;type:text" json:"id"`
	UserID     string  `gorm:"type:text;index;not null" json:"userId"`
	FullName   string  `gorm:"not null" json:"fullName"`
	Phone      string  `gorm:"not null" json:"phone"`
	Street     string  `gorm:"not null" json:"street"`
	City       string  `gorm:"not null" json:"city"`
	State      *string `json:"state,omitempty"`
	PostalCode string  `gorm:"not null" json:"postalCode"`
	Country    string  `gorm:"default:Pakistan" json:"country"`
	IsDefault  bool    `gorm:"default:false" json:"isDefault"`

	User   User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Orders []Order `gorm:"foreignKey:AddressID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Address) TableName() string { return "addresses" }