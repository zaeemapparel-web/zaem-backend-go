package dtos

// ==================== REQUESTS ====================

type AddToWishlistRequest struct {
	ProductID string `json:"productId" binding:"required"`
}

// ==================== RESPONSES ====================

type WishlistItemResponse struct {
	ID            string          `json:"id"`
	ProductID     string          `json:"productId"`
	PriceAtAdd    *float64        `json:"priceAtAdd,omitempty"`
	NotifyOnDrop  bool            `json:"notifyOnDrop"`
	NotifyOnStock bool            `json:"notifyOnStock"`
	Product       ProductResponse `json:"product"`
	CreatedAt     string          `json:"createdAt"`
}