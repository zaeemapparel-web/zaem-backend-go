package services

import (
	"context"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/pkg/cuid"
)

type WishlistService struct {
	repo        repositories.WishlistRepository
	productRepo repositories.ProductRepository
}

func NewWishlistService(r repositories.WishlistRepository, pr repositories.ProductRepository) *WishlistService {
	return &WishlistService{repo: r, productRepo: pr}
}

func (s *WishlistService) GetMyWishlist(ctx context.Context, userID string) ([]dtos.WishlistItemResponse, error) {
	items, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch wishlist")
	}

	resp := make([]dtos.WishlistItemResponse, 0, len(items))
	for i := range items {
		resp = append(resp, s.toResponse(&items[i]))
	}
	return resp, nil
}

func (s *WishlistService) Add(ctx context.Context, userID string, req dtos.AddToWishlistRequest) error {
	// Verify product
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil || product == nil {
		return apperrors.ProductNotFound()
	}

	// Check if already in wishlist
	existing, _ := s.repo.FindByUserAndProduct(ctx, userID, req.ProductID)
	if existing != nil {
		return apperrors.Conflict("Product already in wishlist")
	}

	item := &models.WishlistItem{
		ID:            cuid.New(),
		UserID:        userID,
		ProductID:     req.ProductID,
		PriceAtAdd:    &product.Price,
		NotifyOnDrop:  true,
		NotifyOnStock: true,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperrors.Internal("Failed to add to wishlist")
	}
	return nil
}

func (s *WishlistService) Remove(ctx context.Context, userID, productID string) error {
	existing, _ := s.repo.FindByUserAndProduct(ctx, userID, productID)
	if existing == nil {
		return apperrors.NotFound("Product not in wishlist")
	}

	if err := s.repo.Delete(ctx, userID, productID); err != nil {
		return apperrors.Internal("Failed to remove from wishlist")
	}
	return nil
}

func (s *WishlistService) toResponse(item *models.WishlistItem) dtos.WishlistItemResponse {
	p := item.Product

	prodResp := dtos.ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		Images:       p.Images,
		Stock:        p.Stock,
		IsFeatured:   p.IsFeatured,
	}

	if p.Category.ID != "" {
		prodResp.Category = &dtos.CategoryBrief{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}

	return dtos.WishlistItemResponse{
		ID:            item.ID,
		ProductID:     item.ProductID,
		PriceAtAdd:    item.PriceAtAdd,
		NotifyOnDrop:  item.NotifyOnDrop,
		NotifyOnStock: item.NotifyOnStock,
		Product:       prodResp,
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
	}
}