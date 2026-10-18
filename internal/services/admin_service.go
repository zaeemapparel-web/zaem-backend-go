package services

import (
	"context"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/repositories"
)

type AdminService struct {
	repo repositories.AdminRepository
}

func NewAdminService(repo repositories.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// ==================== DASHBOARD ====================

func (s *AdminService) GetStats(ctx context.Context) (*dtos.DashboardStatsResponse, error) {
	revenue, _ := s.repo.GetTotalRevenue(ctx)
	orders, _ := s.repo.GetTotalOrders(ctx)
	users, _ := s.repo.GetTotalUsers(ctx)
	products, _ := s.repo.GetTotalProducts(ctx)
	pending, _ := s.repo.GetPendingOrders(ctx)
	todayRevenue, _ := s.repo.GetTodayRevenue(ctx)
	todayOrders, _ := s.repo.GetTodayOrders(ctx)
	lowStock, _ := s.repo.GetLowStockCount(ctx, 10)

	return &dtos.DashboardStatsResponse{
		TotalRevenue:  revenue,
		TotalOrders:   orders,
		TotalUsers:    users,
		TotalProducts: products,
		PendingOrders: pending,
		TodayRevenue:  todayRevenue,
		TodayOrders:   todayOrders,
		TotalLowStock: lowStock,
	}, nil
}

func (s *AdminService) GetRecentOrders(ctx context.Context, limit int) ([]dtos.RecentOrderResponse, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	orders, err := s.repo.GetRecentOrders(ctx, limit)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch recent orders")
	}

	result := make([]dtos.RecentOrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, dtos.RecentOrderResponse{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			UserName:    o.User.Name,
			UserEmail:   o.User.Email,
			Total:       o.Total,
			Status:      string(o.Status),
			CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, nil
}

func (s *AdminService) GetTopProducts(ctx context.Context, limit int) ([]dtos.TopProductResponse, error) {
	if limit < 1 || limit > 20 {
		limit = 5
	}

	products, err := s.repo.GetTopProducts(ctx, limit)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch top products")
	}

	result := make([]dtos.TopProductResponse, 0, len(products))
	for _, p := range products {
		img := ""
		if len(p.Images) > 0 {
			img = p.Images[0]
		}
		result = append(result, dtos.TopProductResponse{
			ID:        p.ID,
			Name:      p.Name,
			Image:     img,
			Price:     p.Price,
			SoldCount: p.SoldCount,
			Revenue:   p.Price * float64(p.SoldCount),
		})
	}
	return result, nil
}

// ==================== USER MANAGEMENT ====================

func (s *AdminService) GetAllUsers(ctx context.Context, page, limit int, search string) ([]dtos.AdminUserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	users, total, err := s.repo.GetAllUsers(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch users")
	}

	result := make([]dtos.AdminUserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, dtos.AdminUserResponse{
			ID:            u.ID,
			Name:          u.Name,
			Email:         u.Email,
			Phone:         u.Phone,
			Role:          string(u.Role),
			IsActive:      u.IsActive,
			TotalSpent:    u.TotalSpent,
			LoyaltyPoints: u.LoyaltyPoints,
			Tier:          string(u.Tier),
			CreatedAt:     u.CreatedAt.Format(time.RFC3339),
		})
	}
	return result, total, nil
}

func (s *AdminService) UpdateUserRole(ctx context.Context, adminID, userID, newRole string) error {
	if adminID == userID {
		return apperrors.BadRequest("You cannot change your own role")
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return apperrors.UserNotFound()
	}

	return s.repo.UpdateUserFields(ctx, userID, map[string]interface{}{
		"role": newRole,
	})
}

func (s *AdminService) ToggleUserActive(ctx context.Context, adminID, userID string, isActive bool) error {
	if adminID == userID {
		return apperrors.BadRequest("You cannot deactivate yourself")
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return apperrors.UserNotFound()
	}

	return s.repo.UpdateUserFields(ctx, userID, map[string]interface{}{
		"is_active": isActive,
	})
}