package models

import "time"

type AuditLog struct {
	ID         string  `gorm:"primaryKey;type:text" json:"id"`
	UserID     *string `gorm:"type:text;index" json:"userId,omitempty"`
	Action     string  `gorm:"index;not null" json:"action"` // "user.login", "order.create"
	Entity     string  `gorm:"index" json:"entity,omitempty"` // "user", "order"
	EntityID   *string `gorm:"type:text;index" json:"entityId,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
	UserAgent  *string `json:"userAgent,omitempty"`
	Metadata   *string `gorm:"type:text" json:"metadata,omitempty"` // JSON
	Status     string  `gorm:"default:success" json:"status"`

	CreatedAt time.Time `gorm:"index" json:"createdAt"`
}

func (AuditLog) TableName() string { return "audit_logs" }