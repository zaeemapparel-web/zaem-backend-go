package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerAdminRoutes(r *gin.RouterGroup, h *handlers.AdminHandler) {
	admin := r.Group("/admin")
	admin.Use(middleware.Protect(), middleware.AdminOnly())

	// Dashboard
	admin.GET("/stats", h.GetStats)
	admin.GET("/orders/recent", h.GetRecentOrders)
	admin.GET("/products/top", h.GetTopProducts)
    admin.GET("/orders", h.GetAllOrders)
	// User management
	admin.GET("/users", h.GetUsers)
	admin.PATCH("/users/:id/role", h.UpdateUserRole)
	admin.PATCH("/users/:id/active", h.ToggleUserActive)
}