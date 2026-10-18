package dtos

// ==================== REQUESTS ====================

type AddToCartRequest struct {
	ProductID string  `json:"productId" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1,max=99"`
	Size      *string `json:"size"`
	Color     *string `json:"color"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=99"`
}

// ==================== RESPONSES ====================

// CartResponse — EXACTLY what frontend expects
type CartResponse struct {
	Items     []CartItemResponse `json:"items"`
	ItemCount int                `json:"itemCount"`
	Subtotal  float64            `json:"subtotal"`
}

type CartItemResponse struct {
	ID        string          `json:"id"`
	ProductID string          `json:"productId"`
	Quantity  int             `json:"quantity"`
	Size      *string         `json:"size,omitempty"`
	Color     *string         `json:"color,omitempty"`
	Product   ProductResponse `json:"product"`
}