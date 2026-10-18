package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerProductRoutes(r *gin.RouterGroup, h *handlers.ProductHandler) {
	products := r.Group("/products")

	// Public
	products.GET("", h.List)
	products.GET("/featured", h.GetFeatured)
	products.GET("/:slug", h.GetBySlug)
	products.GET("/:slug/related", h.GetRelated)

	// Admin
	admin := products.Group("")
	admin.Use(middleware.Protect(), middleware.AdminOnly())
	{
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
	}
}