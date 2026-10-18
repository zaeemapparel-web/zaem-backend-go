package repositories

import (
	"context"
	"errors"

	"github.com/zaem/backend-go/internal/database"
	"github.com/zaem/backend-go/internal/models"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository() CartRepository {
	return &cartRepository{db: database.DB}
}

// FindOrCreateByUserID finds or creates a cart for the user
func (r *cartRepository) FindOrCreateByUserID(ctx context.Context, userID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		First(&cart, "user_id = ?", userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new cart
		cart = models.Cart{
			ID:     cuidNew(),
			UserID: userID,
		}
		if err := r.db.WithContext(ctx).Create(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}

	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// Reload reloads the cart with all relations
func (r *cartRepository) Reload(ctx context.Context, userID string) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Category").
		First(&cart, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// FindItemByID finds a cart item by ID
func (r *cartRepository) FindItemByID(ctx context.Context, itemID string) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		First(&item, "id = ?", itemID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// FindItemByProduct finds existing cart item with same product+size+color
func (r *cartRepository) FindItemByProduct(ctx context.Context, cartID, productID string, size, color *string) (*models.CartItem, error) {
	query := r.db.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cartID, productID)

	if size != nil {
		query = query.Where("size = ?", *size)
	} else {
		query = query.Where("size IS NULL")
	}

	if color != nil {
		query = query.Where("color = ?", *color)
	} else {
		query = query.Where("color IS NULL")
	}

	var item models.CartItem
	err := query.First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem adds a new item to cart
func (r *cartRepository) CreateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// UpdateItem updates quantity
func (r *cartRepository) UpdateItem(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// UpdateItemQuantity updates quantity by ID
func (r *cartRepository) UpdateItemQuantity(ctx context.Context, itemID string, qty int) error {
	return r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("id = ?", itemID).
		Update("quantity", qty).
		Error
}

// DeleteItem removes an item
func (r *cartRepository) DeleteItem(ctx context.Context, itemID string) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, "id = ?", itemID).Error
}

// ClearCart removes all items
func (r *cartRepository) ClearCart(ctx context.Context, cartID string) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, "cart_id = ?", cartID).Error
}