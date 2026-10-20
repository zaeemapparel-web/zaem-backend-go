package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerOrderRoutes(r *gin.RouterGroup, h *handlers.OrderHandler) {
	orders := r.Group("/orders")
	orders.Use(middleware.Protect())

	// 1. Admin routes (Must be before /:id)
	admin := orders.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/all", h.AdminGetAllOrders)
		admin.GET("/stats", h.AdminGetStats)
	}
	
	// Admin update status is under /orders/:id/status (not /admin)
	orders.PUT("/:id/status", middleware.AdminOnly(), h.AdminUpdateStatus)

	// 2. User routes
	orders.POST("/", h.Create)
	orders.GET("/", h.GetMyOrders)
	orders.GET("/:id", h.GetByID)
	orders.PUT("/:id/cancel", h.Cancel)
}