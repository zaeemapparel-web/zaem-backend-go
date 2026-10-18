package dtos

// ==================== REQUESTS ====================

type CreateOrderRequest struct {
	AddressID     string  `json:"addressId" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" binding:"required,oneof=COD JAZZCASH EASYPAISA CARD BANK_TRANSFER"`
	CouponCode    *string `json:"couponCode,omitempty"`
	Notes         *string `json:"notes,omitempty"`
	LoyaltyPoints int     `json:"loyaltyPoints" binding:"omitempty,min=0"` // Points to redeem
}

type UpdateOrderStatusRequest struct {
	Status         string  `json:"status" binding:"required,oneof=CONFIRMED PROCESSING SHIPPED DELIVERED CANCELLED RETURNED"`
	TrackingNumber *string `json:"trackingNumber,omitempty"`
	TrackingURL    *string `json:"trackingUrl,omitempty"`
	CancelledReason *string `json:"cancelledReason,omitempty"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required,min=5"`
}

// ==================== RESPONSES ====================

type OrderResponse struct {
	ID            string              `json:"id"`
	OrderNumber   string              `json:"orderNumber"`
	UserID        string              `json:"userId"`
	Subtotal      float64             `json:"subtotal"`
	ShippingCost  float64             `json:"shippingCost"`
	Discount      float64             `json:"discount"`
	Tax           float64             `json:"tax"`
	Total         float64             `json:"total"`
	Status        string              `json:"status"`
	PaymentMethod string              `json:"paymentMethod"`
	PaymentStatus string              `json:"paymentStatus"`
	Notes         *string             `json:"notes,omitempty"`
	CouponCode    *string             `json:"couponCode,omitempty"`
	TrackingNumber *string            `json:"trackingNumber,omitempty"`
	TrackingURL    *string            `json:"trackingUrl,omitempty"`
	Items         []OrderItemResponse `json:"items"`
	Address       *AddressResponse    `json:"address,omitempty"`
	CreatedAt     string              `json:"createdAt"`
	UpdatedAt     string              `json:"updatedAt"`
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Size      *string `json:"size,omitempty"`
	Color     *string `json:"color,omitempty"`
	Image     *string `json:"image,omitempty"`
}