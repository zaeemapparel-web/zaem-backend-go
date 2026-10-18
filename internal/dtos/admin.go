package dtos

// ==================== REQUESTS ====================

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=USER ADMIN MANAGER STAFF"`
}

type ToggleUserActiveRequest struct {
	IsActive bool `json:"isActive"`
}

// ==================== RESPONSES ====================

type DashboardStatsResponse struct {
	TotalRevenue     float64 `json:"totalRevenue"`
	TotalOrders      int64   `json:"totalOrders"`
	TotalUsers       int64   `json:"totalUsers"`
	TotalProducts    int64   `json:"totalProducts"`
	PendingOrders    int64   `json:"pendingOrders"`
	TodayRevenue     float64 `json:"todayRevenue"`
	TodayOrders      int64   `json:"todayOrders"`
	TotalLowStock    int64   `json:"totalLowStock"`
}

type RecentOrderResponse struct {
	ID          string  `json:"id"`
	OrderNumber string  `json:"orderNumber"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	Total       float64 `json:"total"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}

type TopProductResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Price     float64 `json:"price"`
	SoldCount int     `json:"soldCount"`
	Revenue   float64 `json:"revenue"`
}

type AdminUserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      *string `json:"phone,omitempty"`
	Role       string `json:"role"`
	IsActive   bool   `json:"isActive"`
	TotalSpent float64 `json:"totalSpent"`
	LoyaltyPoints int `json:"loyaltyPoints"`
	Tier       string `json:"tier"`
	CreatedAt  string `json:"createdAt"`
}