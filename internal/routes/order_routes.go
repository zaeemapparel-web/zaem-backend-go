package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerOrderRoutes(r *gin.RouterGroup, h *handlers.OrderHandler) {
	orders := r.Group("/orders")
	orders.Use(middleware.Protect())

	// User routes
	orders.POST("", h.Create)
	orders.GET("", h.GetMyOrders)
	orders.GET("/:id", h.GetByID)
	orders.POST("/:id/cancel", h.Cancel)

	// Admin routes
	admin := orders.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		// Frontend ke mutabiq routes
		admin.GET("/stats", h.AdminGetStats)
		admin.GET("/all", h.AdminGetAllOrders)
		admin.PUT("/:id/status", h.AdminUpdateStatus)
	}
}