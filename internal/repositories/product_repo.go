package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{db: database.DB}
}

// ==================== CREATE ====================

func (r *productRepository) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// ==================== READ ====================

func (r *productRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&p, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var p models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&p, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindAll — with filters + pagination
func (r *productRepository) FindAll(ctx context.Context, filter dtos.ProductFilter, pagination dtos.Pagination) ([]models.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Product{})

	// Filter: category slug
	if filter.CategorySlug != "" {
		query = query.Where("category_id IN (SELECT id FROM categories WHERE slug = ?)", filter.CategorySlug)
	}

	// Filter: search
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchTerm)
	}

	// Filter: price range
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	// Filter: featured
	if filter.Featured != nil {
		query = query.Where("is_featured = ?", *filter.Featured)
	}

	// Filter: is active
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	} else {
		// Public default: only active products
		query = query.Where("is_active = ?", true)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort order
	allowed := map[string]string{
		"price":     "price",
		"name":      "name",
		"createdAt": "created_at",
		"popular":   "sold_count",
	}
	orderClause := pagination.OrderClause(allowed, "created_at")

	// Fetch
	var products []models.Product
	err := query.
		Preload("Category").
		Order(orderClause).
		Limit(pagination.Limit).
		Offset(pagination.Offset()).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) FindFeatured(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND is_active = ?", true, true).
		Preload("Category").
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

func (r *productRepository) FindRelated(ctx context.Context, categoryID, excludeID string, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND id != ? AND is_active = ?", categoryID, excludeID, true).
		Limit(limit).
		Find(&products).Error
	return products, err
}

// ==================== UPDATE ====================

func (r *productRepository) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *productRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}

// ==================== DELETE ====================

func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

// ==================== CHECKS ====================

func (r *productRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

// ==================== STATS ====================

func (r *productRepository) IncrementViewCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}

// Helper: none — maybe unused import warning
var _ = fmt.Sprintf