package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type CategoryHandler struct {
	svc *services.CategoryService
}

func NewCategoryHandler(svc *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// GET /api/categories
func (h *CategoryHandler) GetTree(c *gin.Context) {
	result, err := h.svc.GetTree(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"categories": result})
}

// GET /api/categories/:slug
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	result, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"category": result})
}

// POST /api/categories (admin)
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dtos.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, "Category created successfully", gin.H{"category": result})
}

// PUT /api/categories/:id (admin)
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dtos.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Category updated successfully", gin.H{"category": result})
}

// DELETE /api/categories/:id (admin)
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Category deleted successfully", nil)
}