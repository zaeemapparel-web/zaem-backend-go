package services

import (
	"context"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
)

type CartService struct {
	cartRepo    repositories.CartRepository
	productRepo repositories.ProductRepository
}

func NewCartService(cr repositories.CartRepository, pr repositories.ProductRepository) *CartService {
	return &CartService{cartRepo: cr, productRepo: pr}
}

// GetCart
func (s *CartService) GetCart(ctx context.Context, userID string) (*dtos.CartResponse, error) {
	cart, err := s.cartRepo.FindOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch cart")
	}
	return s.toCartResponse(cart), nil
}

// AddItem
func (s *CartService) AddItem(ctx context.Context, userID string, req dtos.AddToCartRequest) (*dtos.CartResponse, error) {
	// Verify product
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil || product == nil {
		return nil, apperrors.ProductNotFound()
	}
	if !product.IsActive {
		return nil, apperrors.ProductNotFound()
	}
	if product.Stock < req.Quantity {
		return nil, apperrors.InsufficientStock(product.Stock)
	}

	// Get or create cart
	cart, err := s.cartRepo.FindOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to add item")
	}

	// Check if same product+size+color exists
	existing, _ := s.cartRepo.FindItemByProduct(ctx, cart.ID, req.ProductID, req.Size, req.Color)

	if existing != nil {
		newQty := existing.Quantity + req.Quantity
		if newQty > product.Stock {
			return nil, apperrors.InsufficientStock(product.Stock)
		}
		if err := s.cartRepo.UpdateItemQuantity(ctx, existing.ID, newQty); err != nil {
			return nil, apperrors.Internal("Failed to add item")
		}
	} else {
		item := &models.CartItem{
			ID:        cuid.New(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Size:      req.Size,
			Color:     req.Color,
		}
		if err := s.cartRepo.CreateItem(ctx, item); err != nil {
			return nil, apperrors.Internal("Failed to add item")
		}
	}

	// Reload
	cart, _ = s.cartRepo.Reload(ctx, userID)
	return s.toCartResponse(cart), nil
}

// UpdateQuantity
func (s *CartService) UpdateQuantity(ctx context.Context, userID, itemID string, req dtos.UpdateCartItemRequest) (*dtos.CartResponse, error) {
	item, err := s.cartRepo.FindItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, apperrors.CartItemNotFound()
	}

	// Verify ownership
	cart, _ := s.cartRepo.FindOrCreateByUserID(ctx, userID)
	if cart == nil || item.CartID != cart.ID {
		return nil, apperrors.CartItemNotFound()
	}

	// Check stock
	if item.Product.Stock < req.Quantity {
		return nil, apperrors.InsufficientStock(item.Product.Stock)
	}

	if err := s.cartRepo.UpdateItemQuantity(ctx, itemID, req.Quantity); err != nil {
		return nil, apperrors.Internal("Failed to update cart")
	}

	cart, _ = s.cartRepo.Reload(ctx, userID)
	return s.toCartResponse(cart), nil
}

// RemoveItem
func (s *CartService) RemoveItem(ctx context.Context, userID, itemID string) (*dtos.CartResponse, error) {
	item, err := s.cartRepo.FindItemByID(ctx, itemID)
	if err != nil || item == nil {
		return nil, apperrors.CartItemNotFound()
	}

	cart, _ := s.cartRepo.FindOrCreateByUserID(ctx, userID)
	if cart == nil || item.CartID != cart.ID {
		return nil, apperrors.CartItemNotFound()
	}

	if err := s.cartRepo.DeleteItem(ctx, itemID); err != nil {
		return nil, apperrors.Internal("Failed to remove item")
	}

	cart, _ = s.cartRepo.Reload(ctx, userID)
	return s.toCartResponse(cart), nil
}

// ClearCart
func (s *CartService) ClearCart(ctx context.Context, userID string) (*dtos.CartResponse, error) {
	cart, err := s.cartRepo.FindOrCreateByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to clear cart")
	}

	if err := s.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return nil, apperrors.Internal("Failed to clear cart")
	}

	cart, _ = s.cartRepo.Reload(ctx, userID)
	return s.toCartResponse(cart), nil
}

// ==================== HELPERS ====================

func (s *CartService) toCartResponse(cart *models.Cart) *dtos.CartResponse {
	resp := &dtos.CartResponse{
		Items:     make([]dtos.CartItemResponse, 0, len(cart.Items)),
		ItemCount: 0,
		Subtotal:  0,
	}

	for i := range cart.Items {
		item := &cart.Items[i]

		// Add item
		itemResp := dtos.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Size:      item.Size,
			Color:     item.Color,
			Product: dtos.ProductResponse{
				ID:           item.Product.ID,
				Name:         item.Product.Name,
				Slug:         item.Product.Slug,
				Price:        item.Product.Price,
				ComparePrice: item.Product.ComparePrice,
				Images:       item.Product.Images,
				Stock:        item.Product.Stock,
			},
		}

		if item.Product.Category.ID != "" {
			itemResp.Product.Category = &dtos.CategoryBrief{
				ID:   item.Product.Category.ID,
				Name: item.Product.Category.Name,
				Slug: item.Product.Category.Slug,
			}
		}

		resp.Items = append(resp.Items, itemResp)
		resp.ItemCount += item.Quantity
		resp.Subtotal += item.Product.Price * float64(item.Quantity)
	}

	return resp
}