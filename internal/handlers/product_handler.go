package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/dtos"
	"github.com/zaem/backend-go/internal/services"
	"github.com/zaem/backend-go/pkg/response"
	customvalidator "github.com/zaem/backend-go/pkg/validator"
)

type ProductHandler struct {
	svc *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// GET /api/products
func (h *ProductHandler) List(c *gin.Context) {
	var filter dtos.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, 400, "Invalid query parameters")
		return
	}

	// Handle featured if present
	if f := c.Query("featured"); f != "" {
		b := f == "true"
		filter.Featured = &b
	}

	pagination := dtos.NewPagination(c.Query("page"), c.Query("limit"), filter.Sort, filter.Order)

	result, total, err := h.svc.List(c.Request.Context(), filter, pagination)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.ListPaginated(c, gin.H{"products": result}, pagination.Page, pagination.Limit, int(total))
}

// GET /api/products/featured
func (h *ProductHandler) GetFeatured(c *gin.Context) {
	limit := 8
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.svc.GetFeatured(c.Request.Context(), limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"products": result})
}

// GET /api/products/:slug
func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	result, err := h.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OK(c, gin.H{"product": result})
}

// GET /api/products/:slug/related
func (h *ProductHandler) GetRelated(c *gin.Context) {
	slug := c.Param("slug")
	limit := 4
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	result, err := h.svc.GetRelated(c.Request.Context(), slug, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.List(c, len(result), gin.H{"products": result})
}

// POST /api/products (admin)
func (h *ProductHandler) Create(c *gin.Context) {
	var req dtos.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, "Product created successfully", gin.H{"product": result})
}

// PUT /api/products/:id (admin)
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dtos.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, customvalidator.FormatErrors(err))
		return
	}

	result, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Product updated successfully", gin.H{"product": result})
}

// DELETE /api/products/:id (admin)
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.OKMessage(c, "Product deleted successfully", nil)
}