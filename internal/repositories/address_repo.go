package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository() AddressRepository {
	return &addressRepository{db: database.DB}
}

func (r *addressRepository) FindByUserID(ctx context.Context, userID string) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) FindByID(ctx context.Context, id string) (*models.Address, error) {
	var addr models.Address
	err := r.db.WithContext(ctx).First(&addr, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *addressRepository) Create(ctx context.Context, addr *models.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *addressRepository) Update(ctx context.Context, addr *models.Address) error {
	return r.db.WithContext(ctx).Save(addr).Error
}

func (r *addressRepository) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("id = ?", id).
		Updates(fields).
		Error
}

func (r *addressRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, "id = ?", id).Error
}

// ClearDefaultForUser sets all addresses' isDefault to false for a user
func (r *addressRepository) ClearDefaultForUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).
		Error
}

func (r *addressRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Address{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}