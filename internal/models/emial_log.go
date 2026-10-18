package models

import "time"

type EmailLog struct {
	ID         string     `gorm:"primaryKey;type:text" json:"id"`
	To         string     `gorm:"not null" json:"to"`
	Subject    string     `gorm:"not null" json:"subject"`
	Template   string     `gorm:"index;not null" json:"template"`
	OrderID    *string    `gorm:"type:text;index" json:"orderId,omitempty"`
	UserID     *string    `gorm:"type:text;index" json:"userId,omitempty"`
	Status     string     `gorm:"index;default:pending" json:"status"`
	Error      *string    `json:"error,omitempty"`
	ResendID   *string    `json:"resendId,omitempty"`
	RetryCount int        `gorm:"default:0" json:"retryCount"`
	SentAt     *time.Time `json:"sentAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (EmailLog) TableName() string { return "email_logs" }