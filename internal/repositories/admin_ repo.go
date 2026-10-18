package repositories

import (
	"context"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
)

type adminRepository struct{}

func NewAdminRepository() AdminRepository {
	return &adminRepository{}
}

// ==================== DASHBOARD STATS ====================

func (r *adminRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	var sum float64
	err := database.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ?", models.OrderDelivered).
		Select("COALESCE(SUM(total), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *adminRepository) GetTotalOrders(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Order{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleUser).Count(&count).Error
	return count, err
}

func (r *adminRepository) GetTotalProducts(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *adminRepository) GetPendingOrders(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&models.Order{}).Where("status = ?", models.OrderPending).Count(&count).Error
	return count, err
}

func (r *adminRepository) GetTodayRevenue(ctx context.Context) (float64, error) {
	var sum float64
	err := database.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("status = ? AND created_at >= CURRENT_DATE", models.OrderDelivered).
		Select("COALESCE(SUM(total), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *adminRepository) GetTodayOrders(ctx context.Context) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("created_at >= CURRENT_DATE").
		Count(&count).Error
	return count, err
}

func (r *adminRepository) GetLowStockCount(ctx context.Context, threshold int) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("stock <= ? AND is_active = ?", threshold, true).
		Count(&count).Error
	return count, err
}

// ==================== RECENT ORDERS ====================

func (r *adminRepository) GetRecentOrders(ctx context.Context, limit int) ([]models.Order, error) {
	var orders []models.Order
	err := database.DB.WithContext(ctx).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// ==================== TOP PRODUCTS ====================

func (r *adminRepository) GetTopProducts(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	err := database.DB.WithContext(ctx).
		Where("sold_count > 0").
		Order("sold_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

// ==================== USER MANAGEMENT ====================

func (r *adminRepository) GetAllUsers(ctx context.Context, limit, offset int, search string) ([]models.User, int64, error) {
	query := database.DB.WithContext(ctx).Model(&models.User{})

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", s, s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, total, err
}

func (r *adminRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := database.DB.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *adminRepository) UpdateUserFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return database.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}