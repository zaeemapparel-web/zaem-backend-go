package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerReviewRoutes(r *gin.RouterGroup, h *handlers.ReviewHandler) {
	// Public (or optional auth)
	r.GET("/products/:productId/reviews", h.GetProductReviews)

	// User (protected)
	reviews := r.Group("/reviews")
	reviews.Use(middleware.Protect())
	{
		reviews.POST("", h.Create)
		reviews.PUT("/:id", h.Update)
		reviews.DELETE("/:id", h.Delete)
		reviews.POST("/:id/helpful", h.MarkHelpful)
	}

	// Admin
	admin := r.Group("/admin/reviews")
	admin.Use(middleware.Protect(), middleware.AdminOnly())
	{
		admin.GET("", h.GetAdminReviews)
		admin.PATCH("/:id", h.AdminAction)
	}
}