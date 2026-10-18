package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerNewsletterRoutes(r *gin.RouterGroup, h *handlers.NewsletterHandler) {
	newsletter := r.Group("/newsletter")
	{
		newsletter.POST("/subscribe", middleware.RateLimitStrict(), h.Subscribe)
		newsletter.POST("/unsubscribe", middleware.RateLimitStrict(), h.Unsubscribe)
	}
}