package dtos

// ==================== REQUESTS ====================

type CreateAddressRequest struct {
	FullName   string  `json:"fullName" binding:"required,min=2,max=100"`
	Phone      string  `json:"phone" binding:"required,min=10,max=15"`
	Street     string  `json:"street" binding:"required,min=5,max=500"`
	City       string  `json:"city" binding:"required,min=2,max=100"`
	State      *string `json:"state"`
	PostalCode string  `json:"postalCode" binding:"required,min=4,max=20"`
	Country    string  `json:"country" binding:"omitempty,max=100"`
	IsDefault  bool    `json:"isDefault"`
}

type UpdateAddressRequest struct {
	FullName   *string `json:"fullName" binding:"omitempty,min=2,max=100"`
	Phone      *string `json:"phone" binding:"omitempty,min=10,max=15"`
	Street     *string `json:"street" binding:"omitempty,min=5,max=500"`
	City       *string `json:"city" binding:"omitempty,min=2,max=100"`
	State      *string `json:"state"`
	PostalCode *string `json:"postalCode" binding:"omitempty,min=4,max=20"`
	Country    *string `json:"country"`
	IsDefault  *bool   `json:"isDefault"`
}

// ==================== RESPONSES ====================

type AddressResponse struct {
	ID         string  `json:"id"`
	FullName   string  `json:"fullName"`
	Phone      string  `json:"phone"`
	Street     string  `json:"street"`
	City       string  `json:"city"`
	State      *string `json:"state,omitempty"`
	PostalCode string  `json:"postalCode"`
	Country    string  `json:"country"`
	IsDefault  bool    `json:"isDefault"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}