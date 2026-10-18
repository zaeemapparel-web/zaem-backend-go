package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerAuthRoutes(r *gin.RouterGroup, h *handlers.AuthHandler) {
	auth := r.Group("/auth")

	// Public routes (with strict rate limit)
	auth.POST("/register", middleware.RateLimitStrict(), h.Register)
	auth.POST("/login", middleware.RateLimitAuth(), h.Login)
	auth.POST("/forgot-password", middleware.RateLimitStrict(), h.ForgotPassword)
	auth.POST("/reset-password", middleware.RateLimitStrict(), h.ResetPassword)

	// Protected routes (need JWT)
	protected := auth.Group("")
	protected.Use(middleware.Protect())
	{
		protected.GET("/me", h.GetMe)
		protected.PUT("/profile", h.UpdateProfile)
		protected.PUT("/change-password", h.ChangePassword)
	}
}