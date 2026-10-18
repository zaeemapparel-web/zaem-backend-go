package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{db: database.DB}
}

// ==================== CREATE ====================

func (r *categoryRepository) Create(ctx context.Context, cat *models.Category) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

// ==================== READ ====================

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	var cat models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		First(&cat, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var cat models.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, name ASC")
		}).
		First(&cat, "slug = ?", slug).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	var cats []models.Category
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, name ASC").
		Find(&cats).Error
	return cats, err
}

// FindTree — top-level categories with 3 levels of children
func (r *categoryRepository) FindTree(ctx context.Context) ([]models.Category, error) {
	var cats []models.Category

	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Order("sort_order ASC, name ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, name ASC").
				Preload("Children", func(db *gorm.DB) *gorm.DB {
					return db.Order("sort_order ASC, name ASC").
						Preload("Children", func(db *gorm.DB) *gorm.DB {
							return db.Order("sort_order ASC, name ASC")
						})
				})
		}).
		Find(&cats).Error

	return cats, err
}

// ==================== UPDATE ====================

func (r *categoryRepository) Update(ctx context.Context, cat *models.Category) error {
	return r.db.WithContext(ctx).Save(cat).Error
}

func (r *categoryRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}

// ==================== DELETE ====================

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Category{}, "id = ?", id).Error
}

// ==================== CHECKS ====================

func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("slug = ?", slug).
		Count(&count).Error
	return count > 0, err
}

func (r *categoryRepository) CountProducts(ctx context.Context, categoryID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error
	return count, err
}

func (r *categoryRepository) CountChildren(ctx context.Context, categoryID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("parent_id = ?", categoryID).
		Count(&count).Error
	return count, err
}