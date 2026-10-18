package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerNotificationRoutes(r *gin.RouterGroup, h *handlers.NotificationHandler) {
	notifs := r.Group("/notifications")
	notifs.Use(middleware.Protect())

	notifs.GET("", h.GetMy)
	notifs.PATCH("/:id/read", h.MarkRead)
	notifs.PATCH("/read-all", h.MarkAllRead)
	notifs.DELETE("/:id", h.Delete)
}