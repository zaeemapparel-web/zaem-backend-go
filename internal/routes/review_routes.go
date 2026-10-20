package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerReviewRoutes(r *gin.RouterGroup, h *handlers.ReviewHandler) {
	// Public
	r.GET("/reviews/product/:productId", h.GetProductReviews)

	// Protected (User)
	reviews := r.Group("/reviews")
	reviews.Use(middleware.Protect())
	{
		reviews.POST("/", h.Create)
		reviews.GET("/my-reviews", h.GetMyReviews)
		reviews.POST("/:id/helpful", h.MarkHelpful)
	}

	// Admin
	admin := r.Group("/reviews/admin")
	admin.Use(middleware.Protect(), middleware.AdminOnly())
	{
		admin.GET("/all", h.GetAdminReviews)
		admin.GET("/stats", h.GetReviewStats)
		admin.PUT("/:id/approve", h.ApproveReview)
		admin.PUT("/:id/reject", h.RejectReview)
		admin.POST("/:id/reply", h.AdminReply)
		admin.DELETE("/:id", h.Delete)
	}
}