package services

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/zaem/backend-go/internal/errors"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/models"
	"github.com/zaem/backend-go/internal/repositories"
	"github.com/zaem/backend-go/internal/utils"
	"github.com/zaem/backend-go/pkg/cuid"
)

type ProductService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(pr repositories.ProductRepository, cr repositories.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  pr,
		categoryRepo: cr,
	}
}

// ==================== PUBLIC ====================

func (s *ProductService) List(ctx context.Context, filter dtos.ProductFilter, pagination dtos.Pagination) ([]dtos.ProductListResponse, int64, error) {
	products, total, err := s.productRepo.FindAll(ctx, filter, pagination)
	if err != nil {
		return nil, 0, apperrors.Internal("Failed to fetch products")
	}

	result := make([]dtos.ProductListResponse, 0, len(products))
	for _, p := range products {
		result = append(result, s.toListResponse(&p))
	}

	return result, total, nil
}

func (s *ProductService) GetBySlug(ctx context.Context, slug string) (*dtos.ProductResponse, error) {
	p, err := s.productRepo.FindBySlug(ctx, strings.ToLower(slug))
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch product")
	}
	if p == nil {
		return nil, apperrors.ProductNotFound()
	}

	// Increment view count (async)
	go func() {
		_ = s.productRepo.IncrementViewCount(context.Background(), p.ID)
	}()

	resp := s.toResponse(p)
	return &resp, nil
}

func (s *ProductService) GetFeatured(ctx context.Context, limit int) ([]dtos.ProductListResponse, error) {
	if limit < 1 || limit > 50 {
		limit = 8
	}

	products, err := s.productRepo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch featured products")
	}

	result := make([]dtos.ProductListResponse, 0, len(products))
	for _, p := range products {
		result = append(result, s.toListResponse(&p))
	}

	return result, nil
}

func (s *ProductService) GetRelated(ctx context.Context, slug string, limit int) ([]dtos.ProductListResponse, error) {
	p, err := s.productRepo.FindBySlug(ctx, strings.ToLower(slug))
	if err != nil || p == nil {
		return nil, apperrors.ProductNotFound()
	}

	if limit < 1 || limit > 20 {
		limit = 4
	}

	products, err := s.productRepo.FindRelated(ctx, p.CategoryID, p.ID, limit)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch related products")
	}

	result := make([]dtos.ProductListResponse, 0, len(products))
	for _, prod := range products {
		result = append(result, s.toListResponse(&prod))
	}

	return result, nil
}

// ==================== ADMIN ====================

func (s *ProductService) Create(ctx context.Context, req dtos.CreateProductRequest) (*dtos.ProductResponse, error) {
	// Verify category exists
	cat, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil || cat == nil {
		return nil, apperrors.CategoryNotFound()
	}

	// Slug
	slug := utils.GenerateSlug(req.Name)
	if req.Slug != "" {
		slug = utils.GenerateSlug(req.Slug)
	}

	exists, err := s.productRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.Internal("Failed to create product")
	}
	if exists {
		return nil, apperrors.ProductSlugExists()
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	p := &models.Product{
		ID:           cuid.New(),
		Name:         strings.TrimSpace(req.Name),
		Slug:         slug,
		Description:  req.Description,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CategoryID:   req.CategoryID,
		Images:       req.Images,
		Sizes:        req.Sizes,
		Colors:       req.Colors,
		Tags:         req.Tags,
		Stock:        req.Stock,
		SKU:          req.SKU,
		IsFeatured:   req.IsFeatured,
		IsActive:     isActive,
	}

	if err := s.productRepo.Create(ctx, p); err != nil {
		return nil, apperrors.Internal("Failed to create product")
	}

	p, _ = s.productRepo.FindByID(ctx, p.ID)
	resp := s.toResponse(p)
	return &resp, nil
}

func (s *ProductService) Update(ctx context.Context, id string, req dtos.UpdateProductRequest) (*dtos.ProductResponse, error) {
	p, err := s.productRepo.FindByID(ctx, id)
	if err != nil || p == nil {
		return nil, apperrors.ProductNotFound()
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Slug != nil {
		newSlug := utils.GenerateSlug(*req.Slug)
		if newSlug != p.Slug {
			exists, _ := s.productRepo.ExistsBySlug(ctx, newSlug)
			if exists {
				return nil, apperrors.ProductSlugExists()
			}
			updates["slug"] = newSlug
		}
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.ComparePrice != nil {
		updates["compare_price"] = *req.ComparePrice
	}
	if req.CategoryID != nil {
		cat, _ := s.categoryRepo.FindByID(ctx, *req.CategoryID)
		if cat == nil {
			return nil, apperrors.CategoryNotFound()
		}
		updates["category_id"] = *req.CategoryID
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}
	if req.Sizes != nil {
		updates["sizes"] = req.Sizes
	}
	if req.Colors != nil {
		updates["colors"] = req.Colors
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.productRepo.UpdateFields(ctx, id, updates); err != nil {
			return nil, apperrors.Internal("Failed to update product")
		}
	}

	p, _ = s.productRepo.FindByID(ctx, id)
	resp := s.toResponse(p)
	return &resp, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	p, err := s.productRepo.FindByID(ctx, id)
	if err != nil || p == nil {
		return apperrors.ProductNotFound()
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal("Failed to delete product")
	}

	return nil
}

// ==================== HELPERS ====================

func (s *ProductService) toResponse(p *models.Product) dtos.ProductResponse {
	resp := dtos.ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Description:  p.Description,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		CategoryID:   p.CategoryID,
		Images:       p.Images,
		Sizes:        p.Sizes,
		Colors:       p.Colors,
		Tags:         p.Tags,
		Stock:        p.Stock,
		SKU:          p.SKU,
		IsFeatured:   p.IsFeatured,
		IsActive:     p.IsActive,
		ViewCount:    p.ViewCount,
		SoldCount:    p.SoldCount,
		CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
	}

	if p.Category.ID != "" {
		resp.Category = &dtos.CategoryBrief{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}

	return resp
}

func (s *ProductService) toListResponse(p *models.Product) dtos.ProductListResponse {
	return dtos.ProductListResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		Images:       p.Images,
		Stock:        p.Stock,
		IsFeatured:   p.IsFeatured,
		CategoryID:   p.CategoryID,
	}
}