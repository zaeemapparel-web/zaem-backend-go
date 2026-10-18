package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerCategoryRoutes(r *gin.RouterGroup, h *handlers.CategoryHandler) {
	categories := r.Group("/categories")

	// Public
	categories.GET("", h.GetTree)
	categories.GET("/:slug", h.GetBySlug)

	// Admin
	admin := categories.Group("")
	admin.Use(middleware.Protect(), middleware.AdminOnly())
	{
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
	}
}