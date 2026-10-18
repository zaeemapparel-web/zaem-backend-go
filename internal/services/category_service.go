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

type CategoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// ==================== PUBLIC ====================

func (s *CategoryService) GetTree(ctx context.Context) ([]dtos.CategoryResponse, error) {
	cats, err := s.repo.FindTree(ctx)
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch categories")
	}

	result := make([]dtos.CategoryResponse, 0, len(cats))
	for i := range cats {
		result = append(result, s.toTreeResponse(&cats[i]))
	}
	return result, nil
}

func (s *CategoryService) GetBySlug(ctx context.Context, slug string) (*dtos.CategoryResponse, error) {
	cat, err := s.repo.FindBySlug(ctx, strings.ToLower(slug))
	if err != nil {
		return nil, apperrors.Internal("Failed to fetch category")
	}
	if cat == nil {
		return nil, apperrors.CategoryNotFound()
	}

	resp := s.toTreeResponse(cat)

	// Product count
	count, _ := s.repo.CountProducts(ctx, cat.ID)
	resp.ProductCount = count

	return &resp, nil
}

// ==================== ADMIN ====================

func (s *CategoryService) Create(ctx context.Context, req dtos.CreateCategoryRequest) (*dtos.CategoryResponse, error) {
	// Generate slug
	slug := utils.GenerateSlug(req.Name)
	if req.Slug != "" {
		slug = utils.GenerateSlug(req.Slug)
	}

	// Check slug exists
	exists, err := s.repo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, apperrors.Internal("Failed to create category")
	}
	if exists {
		return nil, apperrors.CategorySlugExists()
	}

	// Verify parent
	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.repo.FindByID(ctx, *req.ParentID)
		if err != nil || parent == nil {
			return nil, apperrors.NotFound("Parent category not found")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cat := &models.Category{
		ID:          cuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Slug:        slug,
		Description: req.Description,
		Image:       req.Image,
		ParentID:    req.ParentID,
		IsParent:    req.ParentID == nil || *req.ParentID == "",
		SortOrder:   req.SortOrder,
		IsActive:    isActive,
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, apperrors.Internal("Failed to create category")
	}

	// Reload to get parent
	cat, _ = s.repo.FindByID(ctx, cat.ID)
	resp := s.toTreeResponse(cat)
	return &resp, nil
}

func (s *CategoryService) Update(ctx context.Context, id string, req dtos.UpdateCategoryRequest) (*dtos.CategoryResponse, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil || cat == nil {
		return nil, apperrors.CategoryNotFound()
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Slug != nil {
		newSlug := utils.GenerateSlug(*req.Slug)
		if newSlug != cat.Slug {
			exists, _ := s.repo.ExistsBySlug(ctx, newSlug)
			if exists {
				return nil, apperrors.CategorySlugExists()
			}
			updates["slug"] = newSlug
		}
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}
	if req.ParentID != nil {
		if *req.ParentID == "" {
			updates["parent_id"] = nil
			updates["is_parent"] = true
		} else {
			// Prevent self-parenting
			if *req.ParentID == id {
				return nil, apperrors.BadRequest("Category cannot be its own parent")
			}
			updates["parent_id"] = *req.ParentID
			updates["is_parent"] = false
		}
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateFields(ctx, id, updates); err != nil {
			return nil, apperrors.Internal("Failed to update category")
		}
	}

	cat, _ = s.repo.FindByID(ctx, id)
	resp := s.toTreeResponse(cat)
	return &resp, nil
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil || cat == nil {
		return apperrors.CategoryNotFound()
	}

	// Check if products exist
	productCount, _ := s.repo.CountProducts(ctx, id)
	if productCount > 0 {
		return apperrors.CategoryHasProducts()
	}

	// Check if has children
	childCount, _ := s.repo.CountChildren(ctx, id)
	if childCount > 0 {
		return apperrors.BadRequest("Category has sub-categories. Delete them first.")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal("Failed to delete category")
	}

	return nil
}

// ==================== HELPERS ====================

func (s *CategoryService) toTreeResponse(c *models.Category) dtos.CategoryResponse {
	resp := dtos.CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Image:       c.Image,
		ParentID:    c.ParentID,
		IsParent:    c.IsParent,
		SortOrder:   c.SortOrder,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}

	if c.Parent != nil {
		resp.Parent = &dtos.CategoryBrief{
			ID:   c.Parent.ID,
			Name: c.Parent.Name,
			Slug: c.Parent.Slug,
		}
	}

	if len(c.Children) > 0 {
		resp.Children = make([]dtos.CategoryResponse, 0, len(c.Children))
		for i := range c.Children {
			resp.Children = append(resp.Children, s.toTreeResponse(&c.Children[i]))
		}
	}

	return resp
}