package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{db: database.DB}
}

func (r *orderRepository) CreateTx(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Address").
		Preload("User").
		First(&order, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Address").
		First(&order, "order_number = ?", orderNumber).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID string, pagination dtos.Pagination) ([]models.Order, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	err := query.
		Preload("Items").
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset()).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}
func (r *orderRepository) GetAllAdmin(ctx context.Context, limit, offset int, status string) ([]models.Order, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	err := query.
		Preload("User").
		Preload("Items").
		Preload("Address").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}
// ==================== ADMIN STATS ====================

func (r *orderRepository) CountTotal(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Count(&count).Error
	return count, err
}

func (r *orderRepository) CountByStatus(ctx context.Context, status models.OrderStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *orderRepository) SumRevenue(ctx context.Context) (float64, error) {
	var sum float64
	err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("status NOT IN (?)", []models.OrderStatus{models.OrderCancelled, models.OrderReturned}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *orderRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("role = ?", models.RoleUser).Count(&count).Error
	return count, err
}

func (r *orderRepository) CountProducts(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}