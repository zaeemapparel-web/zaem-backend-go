package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zaem/backend-go/internal/handlers"
	"github.com/zaem/backend-go/internal/middleware"
)

func registerOrderRoutes(r *gin.RouterGroup, h *handlers.OrderHandler) {
	orders := r.Group("/orders")
	orders.Use(middleware.Protect()) // All order routes need login

	// User routes
	orders.POST("", h.Create)
	orders.GET("", h.GetMyOrders)
	orders.GET("/:id", h.GetByID)
	orders.POST("/:id/cancel", h.Cancel)

	// Admin routes
	admin := orders.Group("")
	admin.Use(middleware.AdminOnly())
	{
		admin.PATCH("/:id/status", h.AdminUpdateStatus)
	}
}