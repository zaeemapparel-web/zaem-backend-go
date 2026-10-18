package models

import "time"

type Order struct {
	ID            string        `gorm:"primaryKey;type:text" json:"id"`
	OrderNumber   string        `gorm:"uniqueIndex;not null" json:"orderNumber"`
	UserID        string        `gorm:"type:text;index;not null" json:"userId"`
	AddressID     *string       `gorm:"type:text" json:"addressId,omitempty"`
	Subtotal      float64       `gorm:"not null" json:"subtotal"`
	ShippingCost  float64       `gorm:"default:0" json:"shippingCost"`
	Discount      float64       `gorm:"default:0" json:"discount"`
	Tax           float64       `gorm:"default:0" json:"tax"`
	Total         float64       `gorm:"not null" json:"total"`
	Status        OrderStatus   `gorm:"type:text;default:PENDING;index" json:"status"`
	PaymentMethod PaymentMethod `gorm:"type:text;not null" json:"paymentMethod"`
	PaymentStatus PaymentStatus `gorm:"type:text;default:PENDING" json:"paymentStatus"`
	PaymentRef    *string       `json:"paymentRef,omitempty"`
	Notes         *string       `json:"notes,omitempty"`
	CouponID      *string       `gorm:"type:text" json:"couponId,omitempty"`
	CouponCode    *string       `json:"couponCode,omitempty"`
	LoyaltyUsed   int           `gorm:"default:0" json:"loyaltyUsed"`

	// Tracking
	TrackingNumber *string    `json:"trackingNumber,omitempty"`
	TrackingURL    *string    `json:"trackingUrl,omitempty"`
	ShippedAt      *time.Time `json:"shippedAt,omitempty"`
	DeliveredAt    *time.Time `json:"deliveredAt,omitempty"`
	CancelledAt    *time.Time `json:"cancelledAt,omitempty"`

	// Relations
	User    User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Address *Address    `gorm:"foreignKey:AddressID" json:"address,omitempty"`
	Items   []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Coupon  *Coupon     `gorm:"foreignKey:CouponID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID        string  `gorm:"primaryKey;type:text" json:"id"`
	OrderID   string  `gorm:"type:text;index;not null" json:"orderId"`
	ProductID string  `gorm:"type:text;index;not null" json:"productId"`
	Name      string  `gorm:"not null" json:"name"`
	Price     float64 `gorm:"not null" json:"price"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Size      *string `json:"size,omitempty"`
	Color     *string `json:"color,omitempty"`
	Image     *string `json:"image,omitempty"`

	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
}

func (OrderItem) TableName() string {
	return "order_items"
}